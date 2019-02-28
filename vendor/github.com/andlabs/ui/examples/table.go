// 26 august 2018

// +build OMIT

// TODO possible bugs in libui:
// - the checkboxes on macOS retain their values when they shouldn't
// - the table on GTK+ is very thin; the scrolled window needs hexpand=TRUE

package main

import (
	"fmt"
	"image"
	_ "image/png"
	"image/draw"
	"bytes"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

type modelHandler struct {
	row9Text		string
	yellowRow	int
	checkStates	[15]int
}

func newModelHandler() *modelHandler {
	m := new(modelHandler)
	m.row9Text = "You can edit this one"
	m.yellowRow = -1
	return m
}

func (mh *modelHandler) ColumnTypes(m *ui.TableModel) []ui.TableValue {
	return []ui.TableValue{
		ui.TableString(""),		// column 0 text
		ui.TableString(""),		// column 1 text
		ui.TableString(""),		// column 2 text
		ui.TableColor{},			// row background color
		ui.TableColor{},			// column 1 text color
		ui.TableImage{},		// column 1 image
		ui.TableString(""),		// column 4 button text
		ui.TableInt(0),			// column 3 checkbox state
		ui.TableInt(0),			// column 5 progress
	}
}

func (mh *modelHandler) NumRows(m *ui.TableModel) int {
	return 15
}

var img [2]*ui.Image

func (mh *modelHandler) CellValue(m *ui.TableModel, row, column int) ui.TableValue {
	if column == 3 {
		if row == mh.yellowRow {
			return ui.TableColor{1, 1, 0, 1}
		}
		if row == 3 {
			return ui.TableColor{1, 0, 0, 1}
		}
		if row == 11 {
			return ui.TableColor{0, 0.5, 1, 0.5}
		}
		return nil
	}
	if column == 4 {
		if (row % 2) == 1 {
			return ui.TableColor{0.5, 0, 0.75, 1}
		}
		return nil
	}
	if column == 5 {
		if row < 8 {
			return ui.TableImage{img[0]}
		}
		return ui.TableImage{img[1]}
	}
	if column == 7 {
		return ui.TableInt(mh.checkStates[row])
	}
	if column == 8 {
		if row == 0 {
			return ui.TableInt(0)
		}
		if row == 13 {
			return ui.TableInt(100)
		}
		if row == 14 {
			return ui.TableInt(-1)
		}
		return ui.TableInt(50)
	}
	switch column {
	case 0:
		return ui.TableString(fmt.Sprintf("Row %d", row))
	case 2:
		if row == 9 {
			return ui.TableString(mh.row9Text)
		}
		return ui.TableString("Editing this won't change anything")
	case 1:
		return ui.TableString("Colors!")
	case 6:
		return ui.TableString("Make Yellow")
	}
	panic("unreachable")
}

func (mh *modelHandler) SetCellValue(m *ui.TableModel, row, column int, value ui.TableValue) {
	if row == 9 && column == 2 {
		mh.row9Text = string(value.(ui.TableString))
	}
	if column == 6 {		// row background color
		prevYellowRow := mh.yellowRow
		mh.yellowRow = row
		if prevYellowRow != -1 {
			m.RowChanged(prevYellowRow)
		}
		m.RowChanged(mh.yellowRow)
	}
	if column == 7 {		// checkboxes
		mh.checkStates[row] = int(value.(ui.TableInt))
	}
}

func appendImageNamed(i *ui.Image, which string) {
	b := rawImages[which]
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	nr, ok := img.(*image.RGBA)
	if !ok {
		i2 := image.NewRGBA(img.Bounds())
		draw.Draw(i2, i2.Bounds(), img, img.Bounds().Min, draw.Src)
		nr = i2
	}
	i.Append(nr)
}

func setupUI() {
	mainwin := ui.NewWindow("libui Control Gallery", 640, 480, true)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	img[0] = ui.NewImage(16, 16)
	appendImageNamed(img[0], "andlabs_16x16test_24june2016.png")
	appendImageNamed(img[0], "andlabs_32x32test_24june2016.png")
	img[1] = ui.NewImage(16, 16)
	appendImageNamed(img[1], "tango-icon-theme-0.8.90_16x16_x-office-spreadsheet.png")
	appendImageNamed(img[1], "tango-icon-theme-0.8.90_32x32_x-office-spreadsheet.png")

	mh := newModelHandler()
	model := ui.NewTableModel(mh)

	table := ui.NewTable(&ui.TableParams{
		Model:	model,
		RowBackgroundColorModelColumn:	3,
	})
	mainwin.SetChild(table)
	mainwin.SetMargined(true)

	table.AppendTextColumn("Column 1",
		0, ui.TableModelColumnNeverEditable, nil)

	table.AppendImageTextColumn("Column 2",
		5,
		1, ui.TableModelColumnNeverEditable, &ui.TableTextColumnOptionalParams{
			ColorModelColumn:		4,
		});
	table.AppendTextColumn("Editable",
		2, ui.TableModelColumnAlwaysEditable, nil)

	table.AppendCheckboxColumn("Checkboxes",
		7, ui.TableModelColumnAlwaysEditable)
	table.AppendButtonColumn("Buttons",
		6, ui.TableModelColumnAlwaysEditable)

	table.AppendProgressBarColumn("Progress Bar",
		8)

	mainwin.Show()
}

func main() {
	ui.Main(setupUI)
}

var rawImages = map[string][]byte{
	"andlabs_16x16test_24june2016.png": []byte{
  0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
  0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x10,
  0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0xf3, 0xff, 0x61, 0x00, 0x00, 0x00,
  0x01, 0x73, 0x52, 0x47, 0x42, 0x00, 0xae, 0xce, 0x1c, 0xe9, 0x00, 0x00,
  0x00, 0xca, 0x49, 0x44, 0x41, 0x54, 0x38, 0x11, 0xa5, 0x93, 0xb1, 0x0d,
  0xc2, 0x40, 0x0c, 0x45, 0x1d, 0xc4, 0x14, 0x0c, 0x12, 0x41, 0x0f, 0x62,
  0x12, 0x46, 0x80, 0x8a, 0x2e, 0x15, 0x30, 0x02, 0x93, 0x20, 0x68, 0x11,
  0x51, 0x06, 0x61, 0x0d, 0x88, 0x2d, 0x7f, 0xdb, 0x07, 0x87, 0x08, 0xdc,
  0x49, 0x91, 0x7d, 0xf6, 0xf7, 0xf3, 0x4f, 0xa4, 0x54, 0xbb, 0xeb, 0xf6,
  0x41, 0x05, 0x67, 0xcc, 0xb3, 0x9b, 0xfa, 0xf6, 0x17, 0x62, 0xdf, 0xcd,
  0x48, 0x00, 0x32, 0xbd, 0xa8, 0x1d, 0x72, 0xee, 0x3c, 0x47, 0x16, 0xfb,
  0x5c, 0x53, 0x8d, 0x03, 0x30, 0x14, 0x84, 0xf7, 0xd5, 0x89, 0x26, 0xc7,
  0x25, 0x10, 0x36, 0xe4, 0x05, 0xa2, 0x51, 0xbc, 0xc4, 0x1c, 0xc3, 0x1c,
  0xed, 0x30, 0x1c, 0x8f, 0x16, 0x3f, 0x02, 0x78, 0x33, 0x20, 0x06, 0x60,
  0x97, 0x70, 0xaa, 0x45, 0x7f, 0x85, 0x60, 0x5d, 0xb6, 0xf4, 0xc2, 0xc4,
  0x3e, 0x0f, 0x44, 0xcd, 0x1b, 0x20, 0x90, 0x0f, 0xed, 0x85, 0xa8, 0x55,
  0x05, 0x42, 0x43, 0xb4, 0x9e, 0xce, 0x71, 0xb3, 0xe8, 0x0e, 0xb4, 0xc4,
  0xc3, 0x39, 0x21, 0xb7, 0x73, 0xbd, 0xe4, 0x1b, 0xe4, 0x04, 0xb6, 0xaa,
  0x4f, 0x18, 0x2c, 0xee, 0x42, 0x31, 0x01, 0x84, 0xfa, 0xe0, 0xd4, 0x00,
  0xdf, 0xb6, 0x83, 0xf8, 0xea, 0xc2, 0x00, 0x10, 0xfc, 0x1a, 0x05, 0x30,
  0x74, 0x3b, 0xe0, 0xd1, 0x45, 0xb1, 0x83, 0xaa, 0xf4, 0x77, 0x7e, 0x02,
  0x87, 0x1f, 0x42, 0x7f, 0x9e, 0x2b, 0xe8, 0xdf, 0x00, 0x00, 0x00, 0x00,
  0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
},
	"andlabs_32x32test_24june2016.png": []byte{
  0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
  0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x20,
  0x08, 0x06, 0x00, 0x00, 0x00, 0x73, 0x7a, 0x7a, 0xf4, 0x00, 0x00, 0x00,
  0x01, 0x73, 0x52, 0x47, 0x42, 0x00, 0xae, 0xce, 0x1c, 0xe9, 0x00, 0x00,
  0x01, 0x6a, 0x49, 0x44, 0x41, 0x54, 0x58, 0x09, 0xc5, 0x97, 0xc1, 0x0d,
  0xc2, 0x30, 0x0c, 0x45, 0x03, 0x62, 0x0a, 0x06, 0x41, 0x70, 0x07, 0x31,
  0x09, 0x23, 0xc0, 0x89, 0x05, 0x80, 0x11, 0x98, 0x04, 0xc1, 0x15, 0x81,
  0x18, 0x84, 0x35, 0x00, 0x57, 0xfd, 0x8d, 0x13, 0x92, 0x3a, 0x4e, 0x03,
  0x8d, 0x54, 0x39, 0x35, 0xb1, 0xff, 0x8b, 0xed, 0x1e, 0x18, 0xec, 0xae,
  0xdb, 0x97, 0xe9, 0x71, 0x8d, 0x48, 0x7b, 0x33, 0xb9, 0xf5, 0x82, 0xb0,
  0x7f, 0xcc, 0x4c, 0x05, 0x50, 0xa9, 0x2f, 0x26, 0x32, 0xc4, 0xf9, 0x21,
  0x9f, 0xa1, 0x13, 0x8a, 0x5c, 0x16, 0x40, 0x4a, 0x9e, 0x92, 0x14, 0x78,
  0x8a, 0x5c, 0x43, 0xc4, 0xf4, 0x65, 0x3b, 0x01, 0x3c, 0x57, 0x27, 0x43,
  0x4f, 0x97, 0x95, 0x0d, 0x40, 0xc2, 0xe3, 0xe3, 0xb2, 0x7a, 0xba, 0x40,
  0xd8, 0x19, 0x50, 0x5c, 0x03, 0xe2, 0x08, 0x21, 0x10, 0xdf, 0xd7, 0x3a,
  0x88, 0x6c, 0x46, 0xd4, 0x00, 0x5f, 0x42, 0x35, 0x85, 0x03, 0x41, 0x03,
  0xcb, 0x44, 0x00, 0x1a, 0xb2, 0x2e, 0x80, 0x66, 0xd2, 0x43, 0xd9, 0x32,
  0x7c, 0x2e, 0x40, 0x1b, 0x75, 0x0d, 0xe7, 0xdc, 0x94, 0x09, 0xc6, 0x2a,
  0xc3, 0x8e, 0x04, 0xb7, 0x59, 0x43, 0x08, 0x08, 0x64, 0xcc, 0x15, 0xa7,
  0x78, 0x5b, 0x01, 0x45, 0xdf, 0x28, 0x90, 0x43, 0xd0, 0x3e, 0x77, 0x59,
  0x80, 0x8c, 0x0c, 0x5d, 0x84, 0x21, 0xe7, 0x02, 0x94, 0x1c, 0x42, 0x29,
  0x57, 0x3d, 0x6f, 0x16, 0xa0, 0x6d, 0x00, 0x81, 0xcb, 0xec, 0xe1, 0x7e,
  0x61, 0x6f, 0xc6, 0xac, 0xa7, 0x73, 0xfb, 0xae, 0xc8, 0x65, 0x01, 0x6c,
  0xb8, 0xb8, 0x23, 0x71, 0x47, 0xf0, 0x13, 0x11, 0xf2, 0x89, 0x89, 0x3e,
  0x07, 0xd4, 0x5f, 0x41, 0x4c, 0x88, 0x80, 0xfc, 0xaa, 0x14, 0x07, 0x88,
  0x89, 0x43, 0x28, 0x07, 0x42, 0x5d, 0x01, 0x88, 0x95, 0xb2, 0xc9, 0x00,
  0xd2, 0xed, 0x01, 0xa4, 0xad, 0x82, 0x38, 0x84, 0xe8, 0xab, 0x3f, 0x74,
  0x10, 0x0c, 0x59, 0x0e, 0x21, 0xc5, 0xb5, 0x02, 0xa4, 0xde, 0x3a, 0x06,
  0x41, 0x7e, 0x29, 0x47, 0x72, 0x0b, 0x42, 0x22, 0x25, 0x7c, 0x51, 0x00,
  0x89, 0x3c, 0x55, 0x9c, 0xb7, 0x23, 0x14, 0xf3, 0xd5, 0x82, 0x9c, 0x9e,
  0x87, 0x12, 0x73, 0x1f, 0x87, 0xf0, 0x67, 0xc2, 0x01, 0x28, 0x75, 0x6b,
  0x2e, 0x8e, 0x3d, 0x84, 0x7d, 0x8d, 0x68, 0x0b, 0x10, 0xf8, 0x6b, 0xdb,
  0x00, 0xf8, 0x64, 0xbf, 0x12, 0xe6, 0xed, 0x20, 0x8d, 0x0a, 0xe0, 0x5f,
  0xe2, 0xb8, 0x14, 0x87, 0x68, 0x2a, 0x80, 0x1f, 0xff, 0x6d, 0x07, 0x7d,
  0xff, 0x3d, 0x7f, 0x03, 0x93, 0xca, 0x91, 0xa9, 0x89, 0x2a, 0x2e, 0xd2,
  0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
},
	"tango-icon-theme-0.8.90_16x16_x-office-spreadsheet.png": []byte{
  0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
  0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x10,
  0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0xf3, 0xff, 0x61, 0x00, 0x00, 0x00,
  0x06, 0x62, 0x4b, 0x47, 0x44, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0xa0,
  0xbd, 0xa7, 0x93, 0x00, 0x00, 0x00, 0x09, 0x70, 0x48, 0x59, 0x73, 0x00,
  0x00, 0x0b, 0x13, 0x00, 0x00, 0x0b, 0x13, 0x01, 0x00, 0x9a, 0x9c, 0x18,
  0x00, 0x00, 0x00, 0x07, 0x74, 0x49, 0x4d, 0x45, 0x07, 0xd5, 0x04, 0x16,
  0x14, 0x0d, 0x09, 0xd9, 0x88, 0x44, 0xfa, 0x00, 0x00, 0x02, 0x4d, 0x49,
  0x44, 0x41, 0x54, 0x38, 0xcb, 0x95, 0x92, 0xdf, 0x4b, 0x53, 0x61, 0x18,
  0xc7, 0x3f, 0x67, 0x9b, 0x6e, 0x7a, 0xd6, 0xfc, 0x55, 0xe0, 0x59, 0x8a,
  0x95, 0x71, 0x08, 0x52, 0x21, 0x91, 0x2c, 0x8a, 0x0a, 0x94, 0x11, 0x14,
  0x78, 0x21, 0x99, 0x14, 0x81, 0xdd, 0x48, 0x7f, 0x45, 0x63, 0x5d, 0x75,
  0x1b, 0x34, 0xd0, 0x9b, 0x52, 0x83, 0x0a, 0xad, 0x0b, 0x43, 0x72, 0xd0,
  0x85, 0x14, 0x14, 0x68, 0x77, 0x39, 0xcd, 0x6c, 0x94, 0x0c, 0xf4, 0x48,
  0x4d, 0x5b, 0x5b, 0x4b, 0xcf, 0x39, 0x3b, 0xef, 0xe9, 0x62, 0x3f, 0xec,
  0xd7, 0x2e, 0x7c, 0xae, 0xbe, 0x2f, 0xef, 0xf3, 0xfd, 0xbc, 0xdf, 0xf7,
  0x7d, 0x5e, 0x69, 0x78, 0x78, 0xf8, 0xc9, 0xfa, 0xfa, 0x7a, 0x2f, 0xbb,
  0xab, 0xcb, 0xc1, 0x60, 0x70, 0x1c, 0x80, 0x50, 0x28, 0x64, 0xef, 0xb6,
  0x42, 0xa1, 0x90, 0x5d, 0x20, 0xb9, 0x0a, 0x22, 0x12, 0x89, 0xe4, 0x95,
  0x84, 0xa6, 0x69, 0xf8, 0xfd, 0x0a, 0x00, 0x9a, 0xa6, 0xa1, 0x28, 0xfe,
  0xbc, 0x5e, 0x63, 0x60, 0x60, 0xe0, 0x8f, 0x28, 0x45, 0x80, 0xa6, 0x69,
  0x48, 0x52, 0x0e, 0x20, 0x49, 0xb9, 0xf5, 0xce, 0xde, 0x5a, 0xc9, 0xbb,
  0x14, 0x01, 0x8a, 0x5f, 0xc1, 0x2b, 0x7b, 0x01, 0x88, 0xc5, 0x62, 0xf4,
  0xf4, 0xf4, 0x00, 0x10, 0x8d, 0x46, 0x69, 0x69, 0x6d, 0x41, 0xb2, 0x25,
  0xe6, 0xa3, 0xf3, 0xa5, 0x01, 0x86, 0x6e, 0x90, 0x21, 0x83, 0x84, 0x04,
  0x40, 0x2a, 0x95, 0x22, 0x2f, 0x49, 0x7f, 0x4f, 0x91, 0x8f, 0x57, 0x1a,
  0xe0, 0x71, 0xbb, 0x91, 0xbd, 0x5e, 0x0a, 0xaf, 0xe3, 0xf3, 0x55, 0x15,
  0xfc, 0xf8, 0xaa, 0x76, 0x74, 0x49, 0xc0, 0xb6, 0xae, 0xe7, 0x4f, 0xcc,
  0xb5, 0x7e, 0x8c, 0x3c, 0x67, 0xf5, 0xee, 0x1d, 0xb6, 0x16, 0x16, 0x89,
  0xa7, 0x33, 0xb9, 0x26, 0xb9, 0x82, 0xc9, 0xb6, 0x56, 0xbc, 0x6d, 0xed,
  0xff, 0x49, 0xe0, 0xf1, 0x20, 0xcb, 0x32, 0x00, 0xf6, 0xe8, 0x3d, 0x3e,
  0xcd, 0xcd, 0xd1, 0x74, 0xe5, 0x22, 0x65, 0x9d, 0x2a, 0xc2, 0x21, 0x91,
  0x15, 0x02, 0xc3, 0x14, 0x48, 0x0e, 0x0f, 0x4d, 0xcf, 0xa6, 0x78, 0x74,
  0xbc, 0xe3, 0x65, 0xff, 0xec, 0xdb, 0xfe, 0x22, 0x40, 0xd7, 0x75, 0x00,
  0x36, 0x1e, 0x8c, 0x51, 0xb3, 0xb2, 0x4c, 0xc3, 0x8d, 0x3e, 0xb2, 0xe9,
  0x24, 0xc9, 0xf8, 0x2a, 0xa6, 0x10, 0x18, 0x96, 0x8d, 0xb3, 0xaa, 0x16,
  0x53, 0x08, 0xac, 0x4e, 0x15, 0xe7, 0xd2, 0xea, 0x99, 0x11, 0xf5, 0xf0,
  0xed, 0x22, 0xc0, 0xed, 0x76, 0x23, 0xcb, 0x32, 0x89, 0xc8, 0x34, 0xfb,
  0xaf, 0x9e, 0xe7, 0x47, 0x3c, 0x56, 0x34, 0x9a, 0x42, 0xf0, 0xa1, 0xfe,
  0x10, 0x9b, 0x86, 0x0b, 0x53, 0x64, 0x31, 0x24, 0x0b, 0xef, 0xb1, 0x7a,
  0xd4, 0xa7, 0x93, 0x7d, 0xff, 0x24, 0xb0, 0x37, 0xbf, 0x61, 0x19, 0x5b,
  0xe8, 0xd6, 0x8e, 0x59, 0x52, 0xca, 0xd9, 0xd3, 0xe1, 0x26, 0xfc, 0xb5,
  0x0b, 0xc3, 0xd4, 0x51, 0x6b, 0xbd, 0x08, 0xcb, 0xe4, 0xe8, 0xe8, 0x43,
  0x8f, 0xe3, 0xef, 0x04, 0x66, 0x75, 0x0d, 0xce, 0x7d, 0x0d, 0xe0, 0xab,
  0x41, 0x17, 0x16, 0xba, 0x25, 0x70, 0xaa, 0x95, 0x34, 0x37, 0xa7, 0xa9,
  0xcf, 0xae, 0x71, 0xb2, 0x71, 0x2f, 0x5f, 0x52, 0x29, 0x7c, 0x2b, 0xef,
  0xc9, 0x78, 0x3c, 0xdb, 0xc5, 0x04, 0x81, 0x40, 0x20, 0x37, 0x8d, 0xc1,
  0x41, 0xde, 0xdd, 0x1f, 0xe3, 0xe0, 0xf5, 0x5e, 0xea, 0x0e, 0x1c, 0xc1,
  0x76, 0xb9, 0xa8, 0x3e, 0x5d, 0x87, 0xa8, 0x70, 0x10, 0xec, 0x90, 0x78,
  0xfc, 0x39, 0x4d, 0x63, 0x7c, 0x91, 0xf6, 0xc9, 0x09, 0x2a, 0x55, 0x75,
  0x5c, 0x0a, 0x87, 0xc3, 0x53, 0x89, 0x44, 0xe2, 0xc2, 0xef, 0xb3, 0xad,
  0x9d, 0x99, 0x41, 0x7e, 0xf3, 0x1a, 0xf3, 0xdc, 0x09, 0xca, 0x1a, 0xaa,
  0xf1, 0x9c, 0xb5, 0xd1, 0x0d, 0x41, 0x66, 0xc3, 0x64, 0x7a, 0x4a, 0xd0,
  0x33, 0xfb, 0x8a, 0x0a, 0x55, 0x7d, 0x71, 0x6d, 0x61, 0x21, 0x50, 0xea,
  0x7f, 0x30, 0xd2, 0xdd, 0x7d, 0x49, 0x5f, 0x5a, 0xba, 0xe9, 0x4e, 0x26,
  0x55, 0xe7, 0xcf, 0x4c, 0xb9, 0x6d, 0x83, 0x90, 0x65, 0xc3, 0xa5, 0x28,
  0xcb, 0x07, 0xbb, 0xba, 0x6e, 0x9d, 0x1a, 0x1a, 0x9a, 0x00, 0xf8, 0x05,
  0xf4, 0xf5, 0x23, 0xe9, 0x30, 0xeb, 0x2d, 0xf9, 0x00, 0x00, 0x00, 0x00,
  0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
},
	"tango-icon-theme-0.8.90_32x32_x-office-spreadsheet.png": []byte{
  0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
  0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x20,
  0x08, 0x06, 0x00, 0x00, 0x00, 0x73, 0x7a, 0x7a, 0xf4, 0x00, 0x00, 0x00,
  0x04, 0x73, 0x42, 0x49, 0x54, 0x08, 0x08, 0x08, 0x08, 0x7c, 0x08, 0x64,
  0x88, 0x00, 0x00, 0x05, 0xa5, 0x49, 0x44, 0x41, 0x54, 0x58, 0x85, 0xe5,
  0x97, 0xcd, 0x6b, 0x54, 0x57, 0x18, 0xc6, 0x7f, 0xe7, 0xde, 0x3b, 0x93,
  0x51, 0x67, 0x92, 0x8c, 0x66, 0x32, 0x6a, 0x26, 0x4d, 0x94, 0x18, 0x27,
  0xa9, 0x50, 0xbf, 0x5a, 0x2d, 0x94, 0x76, 0xd1, 0x85, 0x05, 0x41, 0x8b,
  0x60, 0x75, 0x53, 0x70, 0x21, 0x0a, 0x55, 0xba, 0x71, 0x51, 0xa8, 0xb4,
  0x7f, 0x40, 0xd7, 0xd5, 0x45, 0xa5, 0x90, 0x6e, 0x5a, 0xa8, 0x20, 0x34,
  0x15, 0x0a, 0x6d, 0xa1, 0x25, 0x58, 0x2c, 0xa4, 0xd6, 0x42, 0x4d, 0xa2,
  0xf9, 0x20, 0xc9, 0xc4, 0x64, 0x3e, 0x4c, 0xa2, 0x13, 0x93, 0xcc, 0xdc,
  0xef, 0x2e, 0x32, 0xf7, 0xf4, 0xde, 0xc9, 0x97, 0x59, 0x74, 0xd5, 0x03,
  0x87, 0x7b, 0xde, 0xf3, 0x9e, 0x73, 0xde, 0xe7, 0x3c, 0xef, 0xf3, 0xde,
  0xb9, 0x03, 0xff, 0xf7, 0x26, 0xaa, 0x27, 0xae, 0x5d, 0xbb, 0x76, 0x4a,
  0xd3, 0xb4, 0x2e, 0x20, 0xa6, 0xaa, 0x2a, 0x8a, 0xa2, 0xe0, 0x38, 0x0e,
  0xb6, 0x6d, 0xcb, 0x6e, 0x59, 0x96, 0x7c, 0xfa, 0xc7, 0x6b, 0xcd, 0x01,
  0xcf, 0x2d, 0xcb, 0x3a, 0xd7, 0xd5, 0xd5, 0x75, 0xcb, 0x1f, 0x4f, 0x5b,
  0x86, 0x48, 0x88, 0x2f, 0xcf, 0x9f, 0x3f, 0x1f, 0xab, 0x8c, 0xe5, 0xbc,
  0xeb, 0xba, 0x81, 0x75, 0x7e, 0xfb, 0x05, 0xc7, 0xb1, 0x4b, 0x97, 0x2e,
  0x7d, 0x09, 0xac, 0x0d, 0xc0, 0xb6, 0xed, 0x7a, 0x80, 0xd1, 0xd1, 0x51,
  0x84, 0x10, 0x12, 0x84, 0xff, 0xe9, 0x07, 0xe6, 0x1f, 0xaf, 0x64, 0x7b,
  0x40, 0x1a, 0x1b, 0x1b, 0x31, 0x0c, 0xa3, 0xbe, 0xda, 0xb7, 0x0c, 0x80,
  0xeb, 0xba, 0x32, 0xc8, 0x8d, 0x1b, 0x37, 0x48, 0x24, 0x12, 0x81, 0xa0,
  0xf9, 0x7c, 0x9e, 0x64, 0x32, 0x29, 0xed, 0x42, 0xa1, 0x40, 0x32, 0x99,
  0x94, 0xfb, 0x0b, 0x85, 0x02, 0x8d, 0x8d, 0x8d, 0xd2, 0x9e, 0x9e, 0x9e,
  0xe6, 0xc2, 0x85, 0x0b, 0xcb, 0x40, 0xad, 0x0a, 0xc0, 0x71, 0x1c, 0x79,
  0x93, 0x54, 0x2a, 0x45, 0x2a, 0x95, 0x0a, 0xdc, 0xaa, 0xa6, 0xa6, 0x26,
  0x30, 0x17, 0x89, 0x44, 0x48, 0xa5, 0x52, 0x72, 0x8f, 0xdf, 0xf6, 0xfc,
  0xde, 0xda, 0xea, 0x34, 0xae, 0x0a, 0xc0, 0x63, 0xc1, 0x3b, 0xb4, 0xfa,
  0x59, 0x9d, 0x9a, 0x95, 0x68, 0xf7, 0xb7, 0xb5, 0xfc, 0x2b, 0x69, 0x40,
  0x8e, 0xa7, 0xa6, 0xa6, 0xb0, 0x6d, 0x3b, 0x10, 0xac, 0x50, 0x28, 0x04,
  0xd6, 0xe4, 0xf3, 0x79, 0x4c, 0xd3, 0x94, 0xfe, 0x7c, 0x3e, 0xef, 0xa9,
  0x5e, 0xfa, 0x37, 0xcc, 0x80, 0x77, 0xd8, 0xce, 0x9d, 0x3b, 0x03, 0xf4,
  0x0a, 0x21, 0xd0, 0x34, 0x8d, 0xe3, 0xc7, 0x8f, 0xcb, 0x43, 0x47, 0x46,
  0x46, 0x68, 0x6b, 0x6b, 0x93, 0x6b, 0x86, 0x87, 0x87, 0xa5, 0x0d, 0x30,
  0x3c, 0x3c, 0x8c, 0xa2, 0x28, 0x2f, 0x0e, 0xc0, 0xbb, 0x9d, 0x10, 0x82,
  0x6c, 0x36, 0x8b, 0xe3, 0x38, 0xcb, 0x44, 0x38, 0x34, 0x34, 0x24, 0xd7,
  0x4f, 0x4c, 0x4c, 0x04, 0x40, 0x67, 0x32, 0x19, 0x69, 0x03, 0x64, 0x32,
  0x19, 0xd2, 0xe9, 0xf4, 0xc6, 0x01, 0x00, 0xec, 0xd8, 0xb1, 0x83, 0xe6,
  0xe6, 0xe6, 0x00, 0x00, 0x55, 0x55, 0xd9, 0xbb, 0x77, 0xaf, 0xb4, 0x43,
  0xa1, 0x50, 0x80, 0x01, 0x4d, 0xd3, 0x02, 0x0c, 0x68, 0x9a, 0xb6, 0xb1,
  0x14, 0xd8, 0xb6, 0x2d, 0x45, 0x98, 0xcb, 0xe5, 0x96, 0x6d, 0x28, 0x14,
  0x0a, 0x0c, 0x0e, 0x0e, 0xca, 0x80, 0xff, 0x19, 0x03, 0x9e, 0x06, 0x9a,
  0x9b, 0x9b, 0xa5, 0x4f, 0x08, 0x41, 0x28, 0x14, 0x92, 0x07, 0x7a, 0xb6,
  0x9f, 0x01, 0x55, 0x55, 0xd9, 0xb3, 0x67, 0x8f, 0xdc, 0xa3, 0xaa, 0xaa,
  0x64, 0xc0, 0x0f, 0x6c, 0x5d, 0x00, 0x80, 0xd4, 0x80, 0x77, 0xb8, 0xa7,
  0xf2, 0x47, 0x8f, 0x1e, 0x49, 0x3b, 0x93, 0xc9, 0x04, 0x40, 0xaf, 0xc4,
  0x40, 0x47, 0x47, 0xc7, 0xc6, 0x18, 0xf0, 0x52, 0xe0, 0x55, 0x81, 0xbf,
  0xf6, 0x55, 0x55, 0x25, 0x9d, 0x4e, 0x07, 0xec, 0xf6, 0xf6, 0x76, 0x09,
  0x40, 0xd3, 0xb4, 0x00, 0x03, 0x7e, 0x0d, 0xbc, 0x10, 0x03, 0xfe, 0x1a,
  0xf6, 0xbf, 0x07, 0xbc, 0x9e, 0xcb, 0xe5, 0x02, 0x1a, 0x18, 0x1f, 0x1b,
  0x23, 0x7b, 0xef, 0x1e, 0xf9, 0xfb, 0xf7, 0x79, 0x36, 0x32, 0xc2, 0x7c,
  0x2e, 0x47, 0x44, 0x51, 0x70, 0x34, 0x8d, 0xf0, 0xd6, 0xad, 0x28, 0x89,
  0x04, 0xe1, 0x73, 0xe7, 0xe8, 0x7c, 0xfb, 0xed, 0x8d, 0xa5, 0xc0, 0xaf,
  0x81, 0xea, 0x2a, 0x48, 0xa7, 0xd3, 0xd8, 0xa6, 0x49, 0xef, 0x57, 0x5f,
  0x31, 0x70, 0xfd, 0x3a, 0xfb, 0x1a, 0x1a, 0x78, 0x6b, 0xd7, 0x2e, 0x5a,
  0x3a, 0x3b, 0x89, 0x1d, 0x3a, 0x84, 0xeb, 0x38, 0x38, 0xa6, 0x49, 0x71,
  0x6e, 0x8e, 0xb1, 0x5c, 0x8e, 0xbe, 0x2b, 0x57, 0xe8, 0x2e, 0x97, 0x89,
  0xd4, 0xd5, 0x71, 0x1a, 0xc2, 0x37, 0xc1, 0x58, 0x17, 0x80, 0xc7, 0xc0,
  0x4a, 0xef, 0x81, 0x3f, 0x7b, 0x7a, 0xb8, 0xf3, 0xc9, 0x27, 0x1c, 0x88,
  0x46, 0xb9, 0x7a, 0xf2, 0x24, 0x9b, 0xc3, 0x61, 0x60, 0x29, 0xc7, 0xb6,
  0x61, 0xe0, 0x98, 0x26, 0x8e, 0x69, 0x12, 0x11, 0x82, 0xf6, 0x86, 0x06,
  0xda, 0xea, 0xea, 0x58, 0x58, 0x5c, 0xe4, 0xa7, 0x81, 0x01, 0xee, 0xc2,
  0x6f, 0x1f, 0xc0, 0xa9, 0xeb, 0x30, 0xb1, 0x2a, 0x00, 0x4f, 0x03, 0x4d,
  0x4d, 0x4d, 0x92, 0x01, 0x8f, 0x15, 0x61, 0x9a, 0xdc, 0xbd, 0x7a, 0x95,
  0x77, 0x3b, 0x3b, 0x39, 0xdc, 0xd6, 0x06, 0x95, 0x94, 0x39, 0xa6, 0x89,
  0x6d, 0x59, 0x32, 0xb8, 0x63, 0x9a, 0xd8, 0xbe, 0x71, 0xc8, 0xb2, 0x78,
  0xa7, 0xb5, 0x95, 0x44, 0x28, 0x74, 0xf8, 0xbb, 0xe1, 0xe1, 0x1f, 0x2f,
  0xc0, 0xd1, 0x2f, 0xa0, 0xb8, 0xae, 0x06, 0x3c, 0x06, 0x3c, 0x00, 0x43,
  0x37, 0x6f, 0x72, 0x28, 0x1e, 0xe7, 0x70, 0x47, 0x07, 0xe8, 0x3a, 0x8e,
  0x6d, 0xaf, 0x18, 0x30, 0x30, 0xe7, 0x03, 0xb6, 0x2f, 0x1e, 0x67, 0x68,
  0xdb, 0xb6, 0x74, 0xff, 0xcc, 0xcc, 0x47, 0xc0, 0xc7, 0x6b, 0xa6, 0xa0,
  0x5a, 0x03, 0x42, 0x08, 0xfa, 0x1f, 0x3e, 0xe4, 0xf5, 0xb3, 0x67, 0x71,
  0x2d, 0x0b, 0xbb, 0x54, 0x5a, 0x37, 0xe0, 0x32, 0x50, 0xae, 0xcb, 0x81,
  0x44, 0x82, 0xbe, 0x99, 0x99, 0xb3, 0xab, 0x02, 0xf0, 0xea, 0x75, 0x6a,
  0x6a, 0x4a, 0x8e, 0x3d, 0x16, 0xf4, 0xe9, 0x69, 0x6a, 0x6b, 0x6b, 0xb1,
  0x9e, 0x3e, 0xc5, 0x29, 0x95, 0xfe, 0x0d, 0x50, 0x09, 0x6a, 0x9b, 0x26,
  0xe3, 0x6f, 0xec, 0xe7, 0xdb, 0xf2, 0x6b, 0xe0, 0xba, 0x98, 0x86, 0x81,
  0x6e, 0x18, 0x58, 0x86, 0x81, 0x6e, 0x98, 0x38, 0xb6, 0x49, 0xab, 0x53,
  0x44, 0x3c, 0xfc, 0xf4, 0x25, 0x40, 0x59, 0x93, 0x01, 0x4f, 0x03, 0x32,
  0xff, 0x42, 0x70, 0xdf, 0x75, 0x79, 0x92, 0xcd, 0xb2, 0x79, 0xcb, 0x16,
  0xec, 0x72, 0x79, 0xd9, 0x0d, 0x6d, 0xd7, 0xa6, 0xfc, 0x7e, 0x03, 0xdb,
  0xba, 0x8b, 0xfc, 0xb9, 0x98, 0x46, 0xd7, 0x0d, 0x0c, 0x43, 0x47, 0xd7,
  0x0d, 0x74, 0xdd, 0x20, 0xa2, 0xb9, 0x44, 0x67, 0x17, 0x01, 0x54, 0x20,
  0xba, 0xaa, 0x06, 0x84, 0x10, 0x4c, 0x4e, 0x4e, 0x06, 0xde, 0x84, 0x00,
  0x4a, 0x2c, 0x46, 0x6f, 0x4f, 0x0f, 0xc9, 0x33, 0x67, 0x50, 0x66, 0x66,
  0xb0, 0xab, 0x58, 0x28, 0xbe, 0x12, 0x27, 0x17, 0xce, 0xb2, 0x6b, 0xff,
  0xdf, 0xf4, 0xf6, 0xec, 0xc2, 0xb2, 0xec, 0x4a, 0xb7, 0x08, 0x6b, 0xd0,
  0x1e, 0xd1, 0xc9, 0x67, 0xc6, 0x88, 0x43, 0x1e, 0x88, 0x28, 0x6b, 0xa5,
  0xa0, 0xa9, 0xa9, 0x89, 0x96, 0x96, 0x16, 0xd9, 0x5b, 0x5b, 0x5b, 0xa9,
  0xdd, 0xb7, 0x8f, 0x81, 0xc9, 0x49, 0xfe, 0xf8, 0xe1, 0x07, 0x4a, 0xc9,
  0x24, 0x6e, 0x32, 0xb9, 0xf4, 0xe9, 0x5d, 0x2e, 0x63, 0x95, 0x4a, 0x64,
  0x0e, 0x0b, 0xa6, 0xf5, 0x31, 0x1e, 0xef, 0x78, 0xcc, 0x2b, 0xa1, 0xaf,
  0x97, 0x2e, 0xe4, 0xda, 0x6c, 0x8f, 0x6f, 0xe2, 0xb5, 0x58, 0x99, 0xc1,
  0xc7, 0x05, 0xe2, 0x7d, 0x3f, 0x53, 0x82, 0x3b, 0x80, 0xbd, 0x6a, 0x0a,
  0x5c, 0xd7, 0x5d, 0x91, 0x01, 0xf7, 0xe0, 0x41, 0x9e, 0xf7, 0xf5, 0x71,
  0x6f, 0x74, 0x94, 0x85, 0x62, 0x91, 0x9d, 0xbb, 0x77, 0xb3, 0xbd, 0xad,
  0x8d, 0x70, 0xa9, 0x84, 0x3d, 0x37, 0x4b, 0xf6, 0x48, 0x84, 0xb2, 0x55,
  0xc4, 0xb2, 0x75, 0xa2, 0x6f, 0x3d, 0xe0, 0xe5, 0xde, 0xd3, 0xd4, 0x1b,
  0x30, 0x3d, 0xd8, 0x47, 0xcf, 0xa2, 0xc6, 0x9e, 0xde, 0x6e, 0xc8, 0x8f,
  0x3c, 0xfb, 0x0b, 0x3e, 0x07, 0x16, 0xd6, 0x2c, 0xc3, 0x8b, 0x17, 0x2f,
  0x06, 0x2a, 0x40, 0x08, 0x81, 0x7d, 0xe2, 0x04, 0x03, 0x6f, 0xbe, 0xc9,
  0xed, 0xcb, 0x97, 0x79, 0x9e, 0xcf, 0xb3, 0xd7, 0x30, 0x98, 0xe8, 0xef,
  0x67, 0x53, 0x34, 0x8a, 0x95, 0x8e, 0x32, 0xf7, 0x57, 0x8c, 0x90, 0x02,
  0xb8, 0xb0, 0x50, 0xd6, 0x31, 0xfe, 0xee, 0xe2, 0xfb, 0x85, 0xa3, 0x34,
  0x17, 0xb3, 0xec, 0xef, 0xbd, 0x85, 0x9e, 0x1b, 0x99, 0x7e, 0x00, 0x1f,
  0x3e, 0x80, 0x7e, 0xa0, 0x5c, 0x0d, 0x20, 0xa6, 0x28, 0xca, 0x9c, 0xeb,
  0xba, 0xb5, 0xf5, 0xf5, 0xff, 0x7e, 0xc2, 0x57, 0x7f, 0x54, 0x1e, 0x3a,
  0x76, 0x8c, 0x86, 0xdb, 0xb7, 0xf9, 0xe5, 0xb3, 0xcf, 0xf8, 0xbd, 0xbb,
  0x9b, 0x58, 0x28, 0x44, 0x52, 0xd7, 0x51, 0x1f, 0x2f, 0x30, 0x3e, 0xfe,
  0x0c, 0x67, 0x93, 0xc0, 0x71, 0x5d, 0x42, 0x33, 0x16, 0x4a, 0x5f, 0x3f,
  0xaf, 0x3e, 0xfc, 0x15, 0xa7, 0x58, 0xa0, 0xb0, 0x75, 0xab, 0x75, 0x07,
  0x2e, 0x66, 0xa1, 0x17, 0x78, 0x02, 0xc1, 0xbf, 0x66, 0x31, 0x20, 0x71,
  0xe4, 0xc8, 0x91, 0xf7, 0x5a, 0x5a, 0x5a, 0x3e, 0x05, 0x36, 0x55, 0xb3,
  0xe3, 0x6f, 0xae, 0xeb, 0x62, 0x18, 0x06, 0xd6, 0xfc, 0x3c, 0x35, 0x93,
  0x93, 0x44, 0x66, 0x67, 0xa9, 0x99, 0x9f, 0x47, 0x33, 0x0c, 0x14, 0xcb,
  0xc2, 0xd1, 0x34, 0xcc, 0x50, 0x08, 0x3d, 0x1a, 0xa5, 0x14, 0x8f, 0x53,
  0xda, 0xbe, 0xdd, 0x18, 0xcd, 0xe5, 0xbe, 0x18, 0x1c, 0x1c, 0xfc, 0x06,
  0x18, 0x07, 0xb2, 0x80, 0xe3, 0x07, 0x10, 0x06, 0x92, 0xc0, 0x36, 0xa0,
  0x1e, 0xa8, 0x61, 0xa9, 0x54, 0xd6, 0x6b, 0x0a, 0x10, 0x01, 0x36, 0x57,
  0xf6, 0x68, 0x95, 0x39, 0x00, 0x87, 0xa5, 0x1f, 0x9e, 0xa7, 0x80, 0x0e,
  0x3c, 0x03, 0x66, 0x59, 0xfa, 0x1d, 0x98, 0xaf, 0x66, 0xc0, 0x6b, 0x35,
  0x95, 0x03, 0xc3, 0xab, 0xf8, 0xd7, 0x6b, 0xa2, 0xd2, 0xbd, 0xaf, 0x0f,
  0xb7, 0x02, 0xc4, 0xac, 0x80, 0xd0, 0x7d, 0x3e, 0xfe, 0x01, 0x9f, 0x1e,
  0x98, 0x64, 0x1e, 0x77, 0xb2, 0x47, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45,
  0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
},
}
