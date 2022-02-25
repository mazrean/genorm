package convert

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/mazrean/genorm/cmd/generator/types"
)

func Convert(tables []*types.Table, joinNum int) ([]*types.Table, []*types.JoinedTable, error) {
	tables, joinedTables, err := convertJoinedTables(tables, joinNum)
	if err != nil {
		return nil, nil, fmt.Errorf("generate joined tables: %w", err)
	}

	return tables, joinedTables, nil
}

type converterTable struct {
	id              int
	table           *types.Table
	refTables       map[int]*converterTable
	refJoinedTables map[int64]*converterJoinedTable
	joinTablesList  []map[int64]*converterJoinedTable
}

type converterJoinedTable struct {
	hash            int64
	tables          map[int]*converterTable
	refTables       map[int]*converterTable
	refJoinedTables map[int64]*converterJoinedTable
}

func newConverterJoinedTable(tables map[int]*converterTable, refTables map[int]*converterTable) *converterJoinedTable {
	gjt := &converterJoinedTable{
		hash:      -1,
		tables:    tables,
		refTables: refTables,
	}

	return gjt
}

func (cjt *converterJoinedTable) tablesHash(tableNum int) int64 {
	if cjt.hash != -1 {
		return cjt.hash
	}

	joinedTableIDs := make([]int, 0, len(cjt.tables))
	for i := range cjt.tables {
		joinedTableIDs = append(joinedTableIDs, i)
	}

	sort.Slice(joinedTableIDs, func(i int, j int) bool {
		return joinedTableIDs[i] < joinedTableIDs[j]
	})

	var joinedTableHash int64 = 0
	var key int64 = 1
	for _, joinTableID := range joinedTableIDs {
		joinedTableHash += int64(joinTableID) * key
		key = key * int64(tableNum) % math.MaxInt64
	}

	cjt.hash = joinedTableHash

	return joinedTableHash
}

func convertJoinedTables(tables []*types.Table, joinNum int) ([]*types.Table, []*types.JoinedTable, error) {
	converterTables, err := tablesToConverterTables(tables, joinNum)
	if err != nil {
		return nil, nil, fmt.Errorf("create generate tables: %w", err)
	}

	converterTables, generateJoinedTableMap := createJoinedTables(converterTables, joinNum)
	converterTables, generateJoinedTableMap = setTablesRefJoinedTable(converterTables, generateJoinedTableMap, joinNum)
	generateJoinedTableMap = setJoinedTablesRefJoinedTable(generateJoinedTableMap, len(tables), joinNum)

	tables, joinedTables, err := converterTableToTable(converterTables, generateJoinedTableMap)
	if err != nil {
		return nil, nil, fmt.Errorf("generate table to table: %w", err)
	}

	return tables, joinedTables, nil
}

func tablesToConverterTables(tables []*types.Table, joinNum int) ([]*converterTable, error) {
	converterTables := make([]*converterTable, 0, len(tables))
	converterTableMap := make(map[string]*converterTable, len(tables))
	for i, table := range tables {
		newconverterTable := &converterTable{
			id:             i,
			table:          table,
			refTables:      map[int]*converterTable{},
			joinTablesList: make([]map[int64]*converterJoinedTable, 0, joinNum-1),
		}

		converterTables = append(converterTables, newconverterTable)
		converterTableMap[table.StructName] = newconverterTable
	}

	for _, converterTableValue := range converterTables {
		for _, refTable := range converterTableValue.table.RefTables {
			refconverterTable, ok := converterTableMap[refTable.Table.StructName]
			if !ok {
				return nil, fmt.Errorf("ref table not found: %s", refTable.Table.StructName)
			}

			converterTableValue.refTables[refconverterTable.id] = refconverterTable
		}

		joinedTable := newConverterJoinedTable(map[int]*converterTable{
			converterTableValue.id: converterTableValue,
		}, converterTableValue.refTables)
		joinedTableMap := map[int64]*converterJoinedTable{
			joinedTable.tablesHash(len(tables)): joinedTable,
		}
		converterTableValue.joinTablesList = append(converterTableValue.joinTablesList, joinedTableMap)
	}

	return converterTables, nil
}

func createJoinedTables(tables []*converterTable, joinNum int) ([]*converterTable, map[int64]*converterJoinedTable) {
	joinedTableHashMap := make(map[int64]*converterJoinedTable)
	for _, table := range tables {
		for _, joinedTable := range table.joinTablesList[0] {
			joinedTableHash := joinedTable.tablesHash(len(tables))
			joinedTableHashMap[joinedTableHash] = joinedTable
		}
	}

	for i := 1; i < joinNum-1; i++ {
		for _, table := range tables {
			joinedTableMap := map[int64]*converterJoinedTable{}
			for _, refconverterTable := range table.refTables {
				for _, refJoinedTable := range refconverterTable.joinTablesList[i-1] {
					// skip if containing the same table
					if _, ok := refJoinedTable.tables[table.id]; ok {
						continue
					}

					joinTables := make(map[int]*converterTable, len(refJoinedTable.tables)+1)
					joinTables[table.id] = table
					for _, table := range refJoinedTable.tables {
						joinTables[table.id] = table
					}

					joinedTableRefs := make(map[int]*converterTable, len(refJoinedTable.refTables)+len(table.refTables)-2)
					for _, refTable := range refJoinedTable.refTables {
						if table.id == refTable.id {
							continue
						}

						joinedTableRefs[refTable.id] = refTable
					}
					for _, refTable := range table.refTables {
						if _, ok := joinedTableRefs[refTable.id]; ok {
							continue
						}
						if _, ok := refJoinedTable.tables[refTable.id]; ok {
							continue
						}

						joinedTableRefs[refTable.id] = refTable
					}

					joinedTable := newConverterJoinedTable(joinTables, joinedTableRefs)
					joinedTableHash := joinedTable.tablesHash(len(tables))

					if _, ok := joinedTableMap[joinedTableHash]; ok {
						continue
					}
					joinedTableMap[joinedTableHash] = joinedTable

					if _, ok := joinedTableHashMap[joinedTableHash]; ok {
						continue
					}
					joinedTableHashMap[joinedTableHash] = joinedTable
				}
			}

			for _, joinedTable := range table.joinTablesList[i-1] {
				for _, refconverterTable := range table.refTables {
					if _, ok := joinedTable.tables[refconverterTable.id]; ok {
						continue
					}

					joinTables := make(map[int]*converterTable, len(joinedTableMap)+1)
					joinTables[refconverterTable.id] = refconverterTable
					for _, table := range joinedTable.tables {
						joinTables[table.id] = table
					}

					joinedTableRefs := map[int]*converterTable{}
					for _, refTable := range joinedTable.refTables {
						if refconverterTable.id == refTable.id {
							continue
						}

						joinedTableRefs[refTable.id] = refTable
					}
					for _, refTable := range refconverterTable.refTables {
						if _, ok := joinedTableRefs[refTable.id]; ok {
							continue
						}
						if _, ok := joinedTable.tables[refTable.id]; ok {
							continue
						}

						joinedTableRefs[refTable.id] = refTable
					}

					joinedTable := newConverterJoinedTable(joinTables, joinedTableRefs)
					joinedTableHash := joinedTable.tablesHash(len(tables))

					if _, ok := joinedTableMap[joinedTableHash]; ok {
						continue
					}
					joinedTableMap[joinedTableHash] = joinedTable

					if _, ok := joinedTableHashMap[joinedTableHash]; ok {
						continue
					}
					joinedTableHashMap[joinedTableHash] = joinedTable
				}
			}

			table.joinTablesList = append(table.joinTablesList, joinedTableMap)
		}
	}

	return tables, joinedTableHashMap
}

func setTablesRefJoinedTable(tables []*converterTable, joinedTableMap map[int64]*converterJoinedTable, joinNum int) ([]*converterTable, map[int64]*converterJoinedTable) {
	for _, table := range tables {
		refJoinedTables := map[int64]*converterJoinedTable{}
		for _, refTable := range table.refTables {
			for _, joinedTables := range refTable.joinTablesList {
				for _, joinedTable := range joinedTables {
					if _, ok := joinedTable.tables[table.id]; ok {
						continue
					}

					joinedTableHash := joinedTable.tablesHash(len(tables))
					refJoinedTables[joinedTableHash] = joinedTable

					joinTables := make(map[int]*converterTable, len(joinedTable.tables)+1)
					joinTables[table.id] = table
					for _, table := range joinedTable.tables {
						joinTables[table.id] = table
					}

					var joinedTableRefs map[int]*converterTable
					if len(joinTables) == joinNum {
						joinedTableRefs = map[int]*converterTable{}
					} else {
						joinedTableRefs = make(map[int]*converterTable, len(joinedTable.refTables)+len(table.refTables)-2)
						for _, refTable := range joinedTable.refTables {
							if table.id == refTable.id {
								continue
							}

							joinedTableRefs[refTable.id] = refTable
						}
						for _, refTable := range table.refTables {
							if _, ok := joinedTableRefs[refTable.id]; ok {
								continue
							}
							if _, ok := joinedTable.tables[refTable.id]; ok {
								continue
							}

							joinedTableRefs[refTable.id] = refTable
						}
					}

					newJoinedTable := newConverterJoinedTable(joinTables, joinedTableRefs)
					newJoinedTableHash := joinedTable.tablesHash(len(tables))

					if _, ok := joinedTableMap[newJoinedTableHash]; ok {
						continue
					}
					joinedTableMap[newJoinedTableHash] = newJoinedTable
				}
			}
		}
		table.refJoinedTables = refJoinedTables
	}

	return tables, joinedTableMap
}

func setJoinedTablesRefJoinedTable(joinedTableMap map[int64]*converterJoinedTable, tableNum int, joinNum int) map[int64]*converterJoinedTable {
	for _, table := range joinedTableMap {
		refJoinedTables := map[int64]*converterJoinedTable{}
		for _, refTable := range table.refTables {
			for _, joinedTables := range refTable.joinTablesList[:joinNum-len(table.tables)] {
			CHECK_LOOP:
				for _, joinedTable := range joinedTables {
					for _, table := range table.tables {
						if _, ok := joinedTable.tables[table.id]; ok {
							continue CHECK_LOOP
						}
					}

					joinedTableHash := joinedTable.tablesHash(tableNum)
					refJoinedTables[joinedTableHash] = joinedTable

					joinTables := make(map[int]*converterTable, len(joinedTable.tables)+len(table.tables))
					for _, table := range table.tables {
						joinTables[table.id] = table
					}
					for _, table := range joinedTable.tables {
						joinTables[table.id] = table
					}

					var joinedTableRefs map[int]*converterTable
					if len(joinTables) == joinNum {
						joinedTableRefs = map[int]*converterTable{}
					} else {
						joinedTableRefs = make(map[int]*converterTable, len(joinedTable.refTables)+len(table.refTables)-2)
						for _, refTable := range joinedTable.refTables {
							if _, ok := table.tables[refTable.id]; ok {
								continue
							}

							joinedTableRefs[refTable.id] = refTable
						}
						for _, refTable := range table.refTables {
							if _, ok := joinedTableRefs[refTable.id]; ok {
								continue
							}
							if _, ok := joinedTable.tables[refTable.id]; ok {
								continue
							}

							joinedTableRefs[refTable.id] = refTable
						}
					}

					newJoinedTable := newConverterJoinedTable(joinTables, joinedTableRefs)
					newJoinedTableHash := joinedTable.tablesHash(tableNum)

					if _, ok := joinedTableMap[newJoinedTableHash]; ok {
						continue
					}
					joinedTableMap[newJoinedTableHash] = newJoinedTable
				}
			}
		}
		table.refJoinedTables = refJoinedTables
	}

	return joinedTableMap
}

func converterTableToTable(converterTables []*converterTable, generateJoinedTableMap map[int64]*converterJoinedTable) ([]*types.Table, []*types.JoinedTable, error) {
	tableMap := make(map[int]*types.Table, len(converterTables))
	for _, converterTable := range converterTables {
		tableMap[converterTable.id] = &types.Table{
			StructName: converterTable.table.StructName,
			Columns:    converterTable.table.Columns,
			Methods:    converterTable.table.Methods,
		}
	}

	joinedTableMap := make(map[int64]*types.JoinedTable, len(generateJoinedTableMap))
	for _, generateJoinedTable := range generateJoinedTableMap {
		if len(generateJoinedTable.tables) < 2 {
			continue
		}

		tables := make([]*types.Table, 0, len(converterTables))
		for _, table := range generateJoinedTable.tables {
			tables = append(tables, tableMap[table.id])
		}

		refTables := make([]*types.RefTable, 0, len(generateJoinedTable.refTables))
		for _, refTable := range generateJoinedTable.refTables {
			table := tableMap[refTable.id]
			refTables = append(refTables, &types.RefTable{
				Table: table,
			})
		}

		joinedTableMap[generateJoinedTable.tablesHash(len(converterTables))] = &types.JoinedTable{
			Tables:    tables,
			RefTables: refTables,
		}
	}

	tables := make([]*types.Table, 0, len(converterTables))
	for _, converterTable := range converterTables {
		table, ok := tableMap[converterTable.id]
		if !ok {
			return nil, nil, errors.New("converterTableToTable: table not found")
		}

		refTables := make([]*types.RefTable, 0, len(converterTable.refTables))
		for _, refTable := range converterTable.refTables {
			table := tableMap[refTable.id]
			refTables = append(refTables, &types.RefTable{
				Table: table,
			})
		}
		table.RefTables = refTables

		refJoinedTables := make([]*types.RefJoinedTable, 0, len(converterTable.refJoinedTables))
		for _, refJoinedTable := range converterTable.refJoinedTables {
			if len(refJoinedTable.tables) < 2 {
				continue
			}

			joinedTable := joinedTableMap[refJoinedTable.tablesHash(len(converterTables))]
			refJoinedTables = append(refJoinedTables, &types.RefJoinedTable{
				Table: joinedTable,
			})
		}
		table.RefJoinedTables = refJoinedTables

		tables = append(tables, table)
	}

	joinedTables := make([]*types.JoinedTable, 0, len(generateJoinedTableMap))
	for _, generateJoinedTable := range generateJoinedTableMap {
		if len(generateJoinedTable.tables) < 2 {
			continue
		}

		joinedTable, ok := joinedTableMap[generateJoinedTable.tablesHash(len(converterTables))]
		if !ok {
			return nil, nil, errors.New("converterTableToTable: joinedTable not found")
		}

		refJoinedTables := make([]*types.RefJoinedTable, 0, len(generateJoinedTable.refJoinedTables))
		for _, refJoinedTable := range generateJoinedTable.refJoinedTables {
			if len(refJoinedTable.tables) < 2 {
				continue
			}

			joinedTable := joinedTableMap[refJoinedTable.tablesHash(len(converterTables))]
			refJoinedTables = append(refJoinedTables, &types.RefJoinedTable{
				Table: joinedTable,
			})
		}
		joinedTable.RefJoinedTables = refJoinedTables

		joinedTables = append(joinedTables, joinedTable)
	}

	return tables, joinedTables, nil
}
