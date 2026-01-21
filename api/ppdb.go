package api

import (
	"fmt"
	"net/http"
	"project_sdu/model"
	"project_sdu/repository"
	"project_sdu/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type PPDBAPI interface {
	Register(c *gin.Context)
	ExportStudentsByBatch(c *gin.Context)
}

type ppdbAPI struct {
	studentService service.StudentService
}

func NewPPDBAPI(studentService service.StudentService) *ppdbAPI {
	return &ppdbAPI{studentService}
}

// ====================
// REGISTER (PUBLIC)
// ====================
func (p *ppdbAPI) Register(c *gin.Context) {
	var student model.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Format data tidak valid",
			Errors: map[string]string{
				"body": "Invalid JSON format",
			},
		})
		return
	}

	// Minimal validation
	errorsMap := make(map[string]string)
	if student.FullName == "" {
		errorsMap["full_name"] = "Nama lengkap wajib diisi"
	}
	if student.Gender == "" {
		errorsMap["gender"] = "Jenis kelamin wajib diisi"
	}
	// Add more as needed...

	if len(errorsMap) > 0 {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Validasi gagal",
			Errors:  errorsMap,
		})
		return
	}

	if err := p.studentService.RegisterPPDB(&student); err != nil {
		switch err {
		case repository.ErrNIKExists:
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Message: "NIK sudah terdaftar",
				Errors:  map[string]string{"nik": "NIK already exists"},
			})
			return

		case repository.ErrNISNExists:
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Message: "NISN sudah terdaftar",
				Errors:  map[string]string{"nisn": "NISN already exists"},
			})
			return
		}

		c.JSON(http.StatusBadRequest, model.ErrorResponse{ 
			Success: false,
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse{
		Success: true,
		Status:  http.StatusCreated,
		Message: "Pendaftaran berhasil! Data Anda telah kami terima.",
		Data:    student,
	})
}

func (p *ppdbAPI) ExportStudentsByBatch(c *gin.Context) {
	batchID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Batch ID tidak valid",
		})
		return
	}

	students, err := p.studentService.GetStudentsByBatchID(batchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	f := excelize.NewFile()
	sheet := "PPDB"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{
		"No", "Nama Lengkap", "NIK", "NISN", "Jenis Kelamin",
		"Tempat Lahir", "Tanggal Lahir", "Agama", "Kewarganegaraan", "Asal Sekolah",

		"Keadaan Ortu", "Status Keluarga", "Anak Ke", "Dari Bersaudara",
		"Tinggal Bersama", "Tinggal Bersama (Lainnya)",

		"Alamat", "RT", "RW", "Desa", "Kecamatan", "Kabupaten", "Provinsi", "Kode Pos",

		"No HP", "Email", "Gol. Darah", "Tinggi (cm)", "Berat (kg)", "Riwayat Penyakit",

		"Nama Ayah", "Pendidikan Ayah", "Pekerjaan Ayah", "Penghasilan Ayah",
		"Nama Ibu", "Pendidikan Ibu", "Pekerjaan Ibu", "Penghasilan Ibu",
		"Nama Wali", "Alamat Ortu/Wali", "No HP Ortu/Wali",

		"Jalur", "Gelombang", "Status Diterima", "Tanggal Daftar",
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, s := range students {
		row := strconv.Itoa(i + 2)

		f.SetCellValue(sheet, "A"+row, i+1)
		f.SetCellValue(sheet, "B"+row, s.FullName)
		f.SetCellValue(sheet, "C"+row, str(s.Nik))
		f.SetCellValue(sheet, "D"+row, str(s.Nisn))
		f.SetCellValue(sheet, "E"+row, s.Gender)
		f.SetCellValue(sheet, "F"+row, str(s.TempatLahir))
		f.SetCellValue(sheet, "G"+row, str(s.TanggalLahir))
		f.SetCellValue(sheet, "H"+row, humanizeEnum(enumPtr(s.Agama)))
		f.SetCellValue(sheet, "I"+row, str(s.Kewarganegaraan))
		f.SetCellValue(sheet, "J"+row, str(s.AsalSekolah))

		f.SetCellValue(sheet, "K"+row, humanizeEnum(enumPtr(s.KeadaanOrtu)))
		f.SetCellValue(sheet, "L"+row, humanizeEnum(enumPtr(s.StatusKeluarga)))
		f.SetCellValue(sheet, "M"+row, intv(s.AnakKe))
		f.SetCellValue(sheet, "N"+row, intv(s.DariBersaudara))
		f.SetCellValue(sheet, "O"+row, humanizeEnum(enumPtr(s.TinggalBersama)))
		f.SetCellValue(sheet, "P"+row, str(s.TinggalBersamaLainnya))

		f.SetCellValue(sheet, "Q"+row, str(s.AlamatJalan))
		f.SetCellValue(sheet, "R"+row, str(s.Rt))
		f.SetCellValue(sheet, "S"+row, str(s.Rw))
		f.SetCellValue(sheet, "T"+row, str(s.DesaKelurahan))
		f.SetCellValue(sheet, "U"+row, str(s.Kecamatan))
		f.SetCellValue(sheet, "V"+row, str(s.Kabupaten))
		f.SetCellValue(sheet, "W"+row, str(s.Provinsi))
		f.SetCellValue(sheet, "X"+row, str(s.KodePos))

		f.SetCellValue(sheet, "Y"+row, str(s.Phone))
		f.SetCellValue(sheet, "Z"+row, str(s.Email))
		f.SetCellValue(sheet, "AA"+row, humanizeEnum(enumPtr(s.BloodType)))
		f.SetCellValue(sheet, "AB"+row, intv(s.TinggiCm))
		f.SetCellValue(sheet, "AC"+row, intv(s.BeratKg))
		f.SetCellValue(sheet, "AD"+row, str(s.RiwayatPenyakit))

		if s.Parent != nil {
			f.SetCellValue(sheet, "AE"+row, str(s.Parent.FatherName))
			f.SetCellValue(sheet, "AF"+row, str(s.Parent.FatherEducation))
			f.SetCellValue(sheet, "AG"+row, str(s.Parent.FatherJob))
			f.SetCellValue(sheet, "AH"+row, str(s.Parent.FatherIncome))
			f.SetCellValue(sheet, "AI"+row, str(s.Parent.MotherName))
			f.SetCellValue(sheet, "AJ"+row, str(s.Parent.MotherEducation))
			f.SetCellValue(sheet, "AK"+row, str(s.Parent.MotherJob))
			f.SetCellValue(sheet, "AL"+row, str(s.Parent.MotherIncome))
			f.SetCellValue(sheet, "AM"+row, str(s.Parent.WaliName))
			f.SetCellValue(sheet, "AN"+row, str(s.Parent.AlamatOrtuWali))
			f.SetCellValue(sheet, "AO"+row, str(s.Parent.NoHpOrtuWali))
		}

		if s.Batch != nil {
			f.SetCellValue(sheet, "AP"+row, s.Batch.Jalur)
			f.SetCellValue(sheet, "AQ"+row, s.Batch.Name)
		}

		f.SetCellValue(sheet, "AR"+row, acceptedText(s.IsAccepted))
		f.SetCellValue(sheet, "AS"+row, s.CreatedAt.Format("2006-01-02"))
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=ppdb_batch_"+strconv.Itoa(batchID)+".xlsx")
	_ = f.Write(c.Writer)
}

func str(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

func intv(v *int) int {
	if v != nil {
		return *v
	}
	return 0
}

// untuk enum POINTER (*Religion, *KeadaanOrtu, dst)
func enumPtr[T any](v *T) string {
	if v != nil {
		return fmt.Sprintf("%v", *v)
	}
	return ""
}

// opsional: biar lebih enak dibaca manusia
func humanizeEnum(v string) string {
	if v == "" {
		return ""
	}
	v = strings.ReplaceAll(v, "_", " ")
	return strings.Title(strings.ToLower(v))
}

func acceptedText(v bool) string {
	if v {
		return "Diterima"
	}
	return "Belum Diterima"
}
