package service

import (
	"context"
	"fmt"

	"github.com/mazrean/gopendb-generator/cmd/usecases/code"
	"github.com/mazrean/gopendb-generator/cmd/usecases/config"
	"github.com/mazrean/gopendb-generator/cmd/usecases/writer"
)

const (
	chanBuf = 10
)

// Generate コード生成のservice
type Generate struct {
	config.Reader
	config.Config
	config.Table
	codeConfig code.Config
	codeTable  code.Table
	writer.Writer
}

// NewGenerate Generateのコンストラクタ
func NewGenerate(cr config.Reader, cf config.Config, ct config.Table, cc code.Config, cdt code.Table, ww writer.Writer) *Generate {
	return &Generate{
		Reader:     cr,
		Config:     cf,
		Table:      ct,
		codeConfig: cc,
		codeTable:  cdt,
		Writer:     ww,
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

	tableDetails, err := g.Table.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get tables: %w", err)
	}

	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	totalChan := make(chan struct{}, chanBuf)
	progressChan := make(chan struct{}, chanBuf)
	progressChans := writer.Progress{
		Total:    totalChan,
		Progress: progressChan,
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
	}(childCtx, totalChan, progressChan)

	isPbOn := true
	err = progressCounter.Start()
	if err != nil {
		isPbOn = false
		cancel()
	}

	for _, tableDetail := range tableDetails {
		tableID := tableDetail.Table.ID

		columns, err := g.Table.GetColumns(tableID)
		if err != nil {
			return fmt.Errorf("failed to get columns: %w", err)
		}

		references, err := g.Table.GetReference(tableID)
		if err != nil {
			return fmt.Errorf("failed to get references: %w", err)
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
			Table:               tableDetail.Table,
			PrimaryKeyColumnIDs: tableDetail.PrimaryKeyColumnIDs,
			Columns:             columns,
			References:          codeReferences,
		}

		fileWriter, err := g.Writer.FileWriterGenerator(ctx, &progressChans, rootPath)
		if err != nil {
			return fmt.Errorf("failed to generate file writer: %w", err)
		}

		g.codeTable.Generate(ctx, &codeTableDetail, fileWriter)
	}

	if isPbOn {
		progressCounter.Finish()
		cancel()
	}

	return nil
}
