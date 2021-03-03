package service

import (
	"context"
	"fmt"

	"github.com/mazrean/gopendb-generator/cmd/usecases/code"
	"github.com/mazrean/gopendb-generator/cmd/usecases/config"
)

const (
	chBuf = 10
)

// Generate コード生成のservice
type Generate struct {
	config.Reader
	config.Config
	config.Table
	codeConfig code.Config
	codeTable  code.Table
}

// NewGenerate Generateのコンストラクタ
func NewGenerate(cr config.Reader, cf config.Config, ct config.Table, cc code.Config, cdt code.Table) *Generate {
	return &Generate{
		Reader:     cr,
		Config:     cf,
		Table:      ct,
		codeConfig: cc,
		codeTable:  cdt,
	}
}

// ProgressCounter コード生成の進捗伝達
type ProgressCounter interface {
	SetTotal(total int)
	Set(progress int)
	Start() error
	Finish() error
}

// Service コード生成のservice
func (g *Generate) Service(ctx context.Context, yamlPath string, rootPath string, progressCounter ProgressCounter) error {
	err := g.Reader.ReadYAML(yamlPath)
	if err != nil {
		return fmt.Errorf("failed to read yaml: %w", err)
	}

	config, err := g.Config.Get()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	g.codeConfig.Set(config)

	tables := g.Table.GetAll()

	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	totalCh := make(chan struct{}, chBuf)
	progressCh := make(chan struct{}, chBuf)
	progressChs := code.Progress{
		Total:    totalCh,
		Progress: progressCh,
	}

	go func(ctx context.Context, totalChan <-chan struct{}, progressChan <-chan struct{}) {
		total := 0
		progress := 0
	GoFuncRoot:
		for {
			select {
			case <-ctx.Done():
				break GoFuncRoot
			case <-totalChan:
				total++
				progressCounter.SetTotal(total)
			case <-progressChan:
				progress++
				progressCounter.Set(progress)
			}
		}
	}(childCtx, totalCh, progressCh)

	isPbOn := true
	err = progressCounter.Start()
	if err != nil {
		isPbOn = false
		cancel()
	}

	for _, table := range tables {
		primaryKeyNames, err := g.Table.GetPrimaryKeyNames(table.ID)
		if err != nil {
			return fmt.Errorf("failed to get primary keys(tableID: %s): %w", table.ID, err)
		}

		columns, err := g.Table.GetColumns(table.ID)
		if err != nil {
			return fmt.Errorf("failed to get columns(tableID: %s): %w", table.ID, err)
		}

		references, err := g.Table.GetReference(table.ID)
		if err != nil {
			return fmt.Errorf("failed to get references(tableID: %s): %w", table.ID, err)
		}

		codeReferences := make([]*code.TableReference, 0, len(references))
		for _, reference := range codeReferences {
			codeReferences = append(codeReferences, &code.TableReference{
				Column:          reference.Column,
				ReferenceTable:  reference.ReferenceTable,
				ReferenceColumn: reference.ReferenceColumn,
			})
		}

		codeTableDetail := code.TableDetail{
			Table:                 table,
			PrimaryKeyColumnNames: primaryKeyNames,
			Columns:               columns,
			References:            codeReferences,
		}

		err = g.codeTable.Generate(ctx, &progressChs, &codeTableDetail)
		if err != nil {
			return fmt.Errorf("failed to generate: %w", err)
		}
	}

	if isPbOn {
		progressCounter.Finish()
		cancel()
	}

	return nil
}
