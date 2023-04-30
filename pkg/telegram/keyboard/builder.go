package keyboard

import (
	tele "gopkg.in/telebot.v3"
)

type Builder struct {
	RowSize int

	rows []tele.Row
}

func NewBuilder(rowSize int) *Builder {
	return &Builder{RowSize: rowSize}
}

func NewBuilderBuffer(rowSize int, rowsCapacity int) *Builder {
	return &Builder{RowSize: rowSize, rows: make([]tele.Row, 0, rowsCapacity)}
}

func (b *Builder) Add(buttons ...tele.Btn) *Builder {
	var row tele.Row
	for i, btn := range buttons {
		row = append(row, btn)
		if (i+1)%b.RowSize == 0 {
			b.rows = append(b.rows, row)
			row = nil
		}
	}

	if len(row) != 0 {
		b.rows = append(b.rows, row)
	}

	return b
}

func (b *Builder) Row(row tele.Row) *Builder {
	copyRow := make(tele.Row, len(row))
	copy(copyRow, row)

	b.rows = append(b.rows, copyRow)
	return b
}

func (b *Builder) OneButtonRow(btn tele.Btn) *Builder {
	b.rows = append(b.rows, tele.Row{btn})
	return b
}

func (b *Builder) Insert(button tele.Btn) *Builder {
	rowsLen := len(b.rows)
	if rowsLen != 0 && len(b.rows[rowsLen-1]) < b.RowSize {
		b.rows[rowsLen-1] = append(b.rows[rowsLen-1], button)
	} else {
		b.Add(button)
	}
	return b
}

// SplitAll splits all buttons on rows with length RowSize or max[0].
func (b *Builder) SplitAll(max ...int) *Builder {
	rowSize := b.RowSize
	if len(max) != 0 {
		rowSize = max[0]
	}
	var total int
	for _, row := range b.rows {
		total += len(row)
	}

	plain := make([]tele.Btn, 0, total)
	for _, row := range b.rows {
		plain = append(plain, row...)
	}
	b.rows = markup.Split(rowSize, plain)

	return b
}

func (b Builder) Inline() *tele.ReplyMarkup {
	m := new(tele.ReplyMarkup)
	m.Inline(b.rows...)
	return m
}

func (b Builder) Reply() *tele.ReplyMarkup {
	m := new(tele.ReplyMarkup)
	m.Reply(b.rows...)
	return m
}
