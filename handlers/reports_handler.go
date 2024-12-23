// handlers/reports_handlers.go
package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"bytes"
    "image"
    "image/jpeg"
    "io/ioutil"
	"math"

	"github.com/gorilla/mux"
	"github.com/jimgustavo/classroom-management/database"
	"github.com/jimgustavo/classroom-management/models"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

//////////////////////////PDF ACTA DE COMPROMISO///////////////////////////

func GenerateReportHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teacherID, err := strconv.Atoi(vars["teacherID"])
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	classroomID, err := strconv.Atoi(vars["classroomID"])
	if err != nil {
		http.Error(w, "Invalid classroom ID", http.StatusBadRequest)
		return
	}

	studentID, err := strconv.Atoi(vars["studentID"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}
	// Fetch classroom
	classroom, err := database.GetClassroomByID(classroomID)
	if err != nil {
		http.Error(w, "Error fetching classroom", http.StatusInternalServerError)
		log.Printf("Error fetching classroom: %v\n", err)
		return
	}
	// Fetch student
	student, err := database.GetStudentByID(studentID)
	if err != nil {
		http.Error(w, "Error fetching student", http.StatusInternalServerError)
		log.Printf("Error fetching student: %v\n", err)
		return
	}
	// Fetch teacher data
	teacherData, err := database.GetTeacherDataByTeacherID(teacherID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving teacher data: %v", err), http.StatusInternalServerError)
		return
	}

	// Generate the PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	// Add a page
	pdf.AddPage()
	// Set font
	pdf.SetFont("Arial", "B", 16)
	// Load template
	err = loadTemplate(pdf, teacherData, classroom, student)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading template: %v", err), http.StatusInternalServerError)
		return
	}
	// Get the current date and time
	currentDatetime := time.Now().Format("2006-01-02_15-04-05")
	// Create the filename in Go format
	filename := fmt.Sprintf("acta-de-compromiso-%s-%s.pdf", student.Name, currentDatetime)

	// Serve the file as a download
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	err = pdf.Output(w)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating PDF: %v", err), http.StatusInternalServerError)
		return
	}
}

func getSpanishMonth(month string) string {
	months := map[string]string{
		"January":   "Enero",
		"February":  "Febrero",
		"March":     "Marzo",
		"April":     "Abril",
		"May":       "Mayo",
		"June":      "Junio",
		"July":      "Julio",
		"August":    "Agosto",
		"September": "Septiembre",
		"October":   "Octubre",
		"November":  "Noviembre",
		"December":  "Diciembre",
	}
	return months[month]
}

func loadTemplate(pdf *gofpdf.Fpdf, teacherData *models.TeacherData, classroom *models.Classroom, student *models.Student) error {

	//In case is needed to use special characters such as: ñ, é, ó and so forth.
	tr := pdf.UnicodeTranslatorFromDescriptor("")
/*
	// Add logo image
	logoPath := "ue12f_logo.jpeg"
	pdf.Image(logoPath, 10, 5, 20, 0, false, "", 0, "ue12f_logo")
*/	 
 // Fetch logo image from URL
 logoURL := fmt.Sprintf("http://localhost:8080/display-logo-as-picture/%d", teacherData.TeacherID)
 resp, err := http.Get(logoURL)
 if err != nil {
	 return fmt.Errorf("failed to fetch logo image: %v", err)
 }
 defer resp.Body.Close()

 if resp.StatusCode != http.StatusOK {
	 return fmt.Errorf("failed to fetch logo image: received status code %d", resp.StatusCode)
 }

 imgData, err := ioutil.ReadAll(resp.Body)
if err != nil {
    return fmt.Errorf("failed to read logo image data: %v", err)
}

img, _, err := image.Decode(bytes.NewReader(imgData))
if err != nil {
    return fmt.Errorf("failed to decode logo image: %v", err)
}

// Get the original dimensions of the image
origWidth := img.Bounds().Dx()
origHeight := img.Bounds().Dy()

// Define the maximum width and height for the logo
const maxWidth = 50.0
const maxHeight = 25.0

// Calculate the scaling factor to maintain the aspect ratio
scaleFactor := math.Min(maxWidth/float64(origWidth), maxHeight/float64(origHeight))

// Calculate the new dimensions
newWidth := float64(origWidth) * scaleFactor
newHeight := float64(origHeight) * scaleFactor

// Convert image to JPEG format
var jpegBuffer bytes.Buffer
err = jpeg.Encode(&jpegBuffer, img, nil)
if err != nil {
    return fmt.Errorf("failed to encode logo image to JPEG: %v", err)
}

// Add logo image to PDF with restricted size
pdf.RegisterImageOptionsReader("logo", gofpdf.ImageOptions{ImageType: "JPEG"}, bytes.NewReader(jpegBuffer.Bytes()))
pdf.ImageOptions("logo", 10, 5, newWidth, newHeight, false, gofpdf.ImageOptions{ImageType: "JPEG"}, 0, "")

	// Your template content
	schoolName := teacherData.School
	studentName := student.Name
	classroomName := classroom.Name
	//schoolYear := teacherData.SchoolYear
	//countryCityDate := fmt.Sprintf("%s, %s %s", teacherData.Country, teacherData.City, time.Now().Format("02 January 2006"))

	// Get the current date components
	day := time.Now().Format("02")
	month := getSpanishMonth(time.Now().Format("January"))
	year := time.Now().Format("2006")
	// Adding content to the document
	// Add title
	pdf.CellFormat(180, 10, schoolName, "0", 0, "C", false, 0, "")
	pdf.Ln(10)
	pdf.Cell(40, 30, tr("ACTA DE COMPROMISO POR BAJO RENDIMIENTO ACADéMICO"))
	pdf.Ln(20)

	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(0, 10, fmt.Sprintf(tr("En la ciudad de %s, a los %s días del mes de %s del %s, comparecen ante el/la rector(a)/vicerrector(a)/tutor(a) de la %s, el/la Sr./Sra. ____________________________, en calidad de representante legal del estudiante %s del curso %s, para suscribir la presente acta de compromiso."),
		teacherData.City, day, month, year, teacherData.School, studentName, classroomName), "", "", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 10, tr("Considerando:"), "", "", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 10, tr("1. Que, conforme al Artículo 26 de la Constitución de la República del Ecuador, \"la educación es un derecho de las personas a lo largo de su vida y un deber ineludible e inexcusable del Estado. Constituye un área prioritaria de la política pública y de la inversión estatal, garantía de la igualdad e inclusión social y condición indispensable para el buen vivir. Las personas, las familias y la sociedad tienen el derecho y la responsabilidad de participar en el proceso educativo.\""), "", "", false)
	pdf.MultiCell(0, 10, tr("2. Que, de acuerdo con el Artículo 8 de la Ley Orgánica de Educación Intercultural (LOEI), las y los estudiantes tienen obligaciones y responsabilidades, tales como cumplir con las actividades académico-formativas, participar en evaluaciones, procurar la excelencia educativa y mostrar integridad y honestidad académica."), "", "", false)
	pdf.MultiCell(0, 10, tr("3. Que, según el Artículo 13 de la LOEI, las madres, padres y/o representantes de los estudiantes deben involucrarse activamente en los procesos educativos de sus representados y atender los llamados y requerimientos de los profesores y autoridades de los planteles, y apoyar y motivar a sus representados especialmente cuando existan dificultades en el proceso de aprendizaje."), "", "", false)
	pdf.MultiCell(0, 10, tr("4. Que, el Artículo 32 del Reglamento a la Ley Orgánica de Educación Intercultural (RLOEI) establece que, si la evaluación continua determinare bajos resultados en los procesos de aprendizaje, se deberá diseñar e implementar de inmediato procesos de refuerzo pedagógico."), "", "", false)
	pdf.Ln(5)

	pdf.MultiCell(0, 10, tr("Compromisos:"), "", "", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 10, fmt.Sprintf(tr("El/la Sr./Sra. ________________________________, en su calidad de representante legal del estudiante %s, se compromete a:"), studentName), "", "", false)
	pdf.MultiCell(0, 10, fmt.Sprintf(tr("1. Garantizar la asistencia regular del estudiante %s a todas las clases y actividades educativas programadas."), studentName), "", "", false)
	pdf.MultiCell(0, 10, tr("2. Colaborar activamente con el personal docente en el diseño e implementación de un plan de refuerzo pedagógico, que incluirá clases de refuerzo, tutorías individuales y un cronograma de estudios a seguir en casa."), "", "", false)
	pdf.MultiCell(0, 10, tr("3. Proveer un ambiente de aprendizaje adecuado en el hogar, dedicando espacios específicos para las tareas y estudios del estudiante."), "", "", false)
	pdf.MultiCell(0, 10, tr("4. Motivar y apoyar al estudiante en el cumplimiento de sus obligaciones académicas, fomentando la excelencia educativa y la integridad académica."), "", "", false)
	pdf.MultiCell(0, 10, tr("5. Participar en reuniones de seguimiento con los docentes y autoridades del plantel para evaluar el progreso académico del estudiante."), "", "", false)
	pdf.MultiCell(0, 10, tr("6. Reconocer y valorar los esfuerzos y avances del estudiante, así como los méritos y la excelencia del personal docente."), "", "", false)
	pdf.Ln(5)

	// Calificaciones actuales del estudiante
	pdf.MultiCell(0, 10, tr("Calificaciones actuales del estudiante:"), "", "", false)
	pdf.Ln(5)

	// Fetch grades below seven
	gradesData, err := database.FetchGradesBelowSevenByClassroomID(classroom.ID)
	if err != nil {
		return fmt.Errorf("error fetching grades below seven: %w", err)
	}

	// Iterate over the fetched grades data for the target student
	for _, studentTermGrades := range gradesData.Grades {
		if studentTermGrades.StudentID == student.ID {
			// Fetch the subject name using the subject ID
			subject, err := database.GetSubjectByID(studentTermGrades.SubjectID)
			if err != nil {
				return fmt.Errorf("error fetching subject: %w", err)
			}
			// Display Student ID
			//pdf.MultiCell(0, 10, fmt.Sprintf(tr("Estudiante ID: %d"), studentTermGrades.StudentID), "", "", false)
			//pdf.Ln(5)
			for _, termGradeSkills := range studentTermGrades.Terms {
				// Display term and subject ID
				pdf.MultiCell(0, 5, fmt.Sprintf(tr("Periodo: %s, Asignatura: %s"), termGradeSkills.Term, subject.Name), "", "", false)

				// Iterate over each skill and grade
				for _, gradeSkill := range termGradeSkills.Grades {
					date := gradeSkill.Date[:10] // Format the date to remove the time part
					pdf.MultiCell(0, 5, fmt.Sprintf(tr("Destreza: %s, Fecha: %s, Nota: %.1f"), gradeSkill.Skill, date, gradeSkill.Grade), "", "", false)
				}
				pdf.Ln(5)
			}
		}
	}

	// Signatures Section
	pdf.MultiCell(0, 10, tr("Firman en conformidad:"), "", "", false)
	pdf.Ln(10)

	// Row 1: Representante Legal and Rector(a)/Vicerrector(a)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(20, 10, "")
	pdf.CellFormat(40, 10, tr("Representante Legal"), "0", 0, "L", false, 0, "")
	pdf.Cell(70, 10, "")
	pdf.CellFormat(20, 10, tr("Rector(a)/Vicerrector(a)"), "0", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 9)
	pdf.Ln(5)
	pdf.Cell(20, 10, "")
	pdf.CellFormat(40, 10, "CI:", "0", 0, "L", false, 0, "")
	pdf.Cell(70, 10, "")
	pdf.CellFormat(40, 10, "CI:", "0", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.Cell(20, 10, "")
	pdf.CellFormat(40, 10, "Correo:", "0", 0, "L", false, 0, "")
	pdf.Cell(70, 10, "")
	pdf.CellFormat(40, 10, "Correo:", "0", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.Cell(20, 10, "")
	pdf.CellFormat(40, 10, tr("Teléfono:"), "0", 0, "L", false, 0, "")
	pdf.Cell(70, 10, "")
	pdf.CellFormat(40, 10, tr("Teléfono:"), "0", 0, "L", false, 0, "")
	pdf.Ln(20)

	// Row 2: Estudiante and Docente Tutor/a
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(20, 10, "")
	pdf.CellFormat(40, 10, tr("Estudiante"), "0", 0, "L", false, 0, "")
	pdf.Cell(70, 10, "")
	pdf.CellFormat(20, 10, tr("Docente Tutor/a"), "0", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 9)
	pdf.Ln(5)
	pdf.Cell(20, 10, "")
	pdf.CellFormat(40, 10, "CI:", "0", 0, "L", false, 0, "")
	pdf.Cell(70, 10, "")
	pdf.CellFormat(20, 10, fmt.Sprintf("CI: %s", teacherData.TeacherIDNumber), "0", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.Cell(20, 10, "")
	pdf.CellFormat(40, 10, "Correo:", "0", 0, "L", false, 0, "")
	pdf.Cell(70, 10, "")
	pdf.CellFormat(20, 10, fmt.Sprintf("Correo: %s", teacherData.InstitutionalEmail), "0", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.Cell(20, 10, "")
	pdf.CellFormat(40, 10, tr("Teléfono:"), "0", 0, "L", false, 0, "")
	pdf.Cell(70, 10, "")
	pdf.CellFormat(20, 10, fmt.Sprintf(tr("Teléfono: %s"), teacherData.Phone), "0", 0, "L", false, 0, "")

	return nil
}

// //////////////////////////////XLSX AVERAGE REPORT///////////////////////////
func GenerateFinalAveragesReport(w http.ResponseWriter, r *http.Request) {
	// Extract the classroom ID and academic period ID from the URL path
	vars := mux.Vars(r)
	classroomIDStr := vars["classroomID"]
	academicPeriodIDStr := vars["academicPeriodID"]
	teacherIDStr := vars["teacherID"]

	// Convert classroomID and academicPeriodID to integers
	classroomID, err := strconv.Atoi(classroomIDStr)
	if err != nil {
		http.Error(w, "Invalid classroom ID", http.StatusBadRequest)
		return
	}

	academicPeriodID, err := strconv.Atoi(academicPeriodIDStr)
	if err != nil {
		http.Error(w, "Invalid academic period ID", http.StatusBadRequest)
		return
	}

	teacherID, err := strconv.Atoi(teacherIDStr)
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	// Fetch students, subjects, and averages with reinforcement grades
	classroom, err := database.GetClassroomByID(classroomID)
	if err != nil {
		http.Error(w, "Error fetching the classroom", http.StatusInternalServerError)
		log.Printf("Error fetching the classroom: %v\n", err)
		return
	}

	students, err := database.GetStudentsByClassroomID(classroomID)
	if err != nil {
		http.Error(w, "Error fetching students", http.StatusInternalServerError)
		log.Printf("Error fetching students: %v\n", err)
		return
	}

	subjects, err := database.GetSubjectsInClassroom(classroomID)
	if err != nil {
		http.Error(w, "Error fetching subjects", http.StatusInternalServerError)
		log.Printf("Error fetching subjects: %v\n", err)
		return
	}

	// Hardcoded term factors
	termFactors := []models.TermFactor{
		{Term: "bimestre_1", Factor: 0.7},
		{Term: "sumativa_1", Factor: 0.3},
		{Term: "bimestre_2", Factor: 0.7},
		{Term: "sumativa_2", Factor: 0.3},
	}

	averagesData, err := database.FetchAveragesWithReinforcementByClassroomID(classroomID, termFactors)
	if err != nil {
		http.Error(w, "Error fetching averages", http.StatusInternalServerError)
		log.Printf("Error fetching averages: %v\n", err)
		return
	}

	teacherData, err := database.GetTeacherDataByTeacherID(teacherID)
	if err != nil {
		http.Error(w, "Error fetching teacher data", http.StatusInternalServerError)
		log.Printf("Error fetching teacher data: %v\n", err)
		return
	}

	terms, err := database.FetchTermsByAcademicPeriodFromDB(academicPeriodID)
	if err != nil {
		http.Error(w, "Error fetching terms", http.StatusInternalServerError)
		log.Printf("Error fetching terms: %v\n", err)
		return
	}

	file := excelize.NewFile()

	for _, subject := range subjects {
		sheetName := subject.Name
		file.NewSheet(sheetName)

		currentTime := time.Now()
		truncatedTime := currentTime.Truncate(time.Second)

		headerSheetStyle, err := file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "center",
			},
			Font: &excelize.Font{
				Bold: true,
			},
		})
		if err != nil {
			fmt.Println(err)
			return
		}

		headerStyle, err := file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Vertical:     "center",
				ShrinkToFit:  true,
				TextRotation: 90,
			},
			Font: &excelize.Font{
				Bold: true,
			},
		})
		if err != nil {
			fmt.Println(err)
			return
		}

		otherHeaderStyle, err := file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Vertical:    "center",
				ShrinkToFit: false,
			},
			Font: &excelize.Font{
				Bold: true,
			},
		})

		if err != nil {
			fmt.Println(err)
			return
		}

		file.SetCellValue(subject.Name, "B1", teacherData.School)
		file.SetCellStyle(sheetName, "B1", "B1", headerSheetStyle)
		file.SetCellValue(subject.Name, "B2", teacherData.Country+" - "+teacherData.City)
		file.SetCellStyle(sheetName, "B2", "B2", headerSheetStyle)
		file.SetCellValue(subject.Name, "B3", "Teacher: "+teacherData.TeacherFullName)
		file.SetCellValue(subject.Name, "B4", "Classroom: "+classroom.Name)
		file.SetCellValue(subject.Name, "D3", "Subject: "+subject.Name)
		file.SetCellValue(subject.Name, "J3", "School Year: "+teacherData.SchoolYear)
		file.SetCellValue(subject.Name, "D4", "Report Date: "+truncatedTime.Format("2006-01-02 15:04:05"))
		file.SetCellValue(subject.Name, "J4", "School Hours: "+teacherData.SchoolHours)
		file.SetCellValue(subject.Name, "B5", "Promedio Anual")
		file.SetCellStyle(sheetName, "B5", "B5", headerSheetStyle)

		mergeCellRanges := [][]string{{"B1", "O1"}, {"B2", "O2"}, {"D3", "I3"}, {"J3", "O3"}, {"D4", "I4"}, {"J4", "O4"}, {"B5", "O5"}}
		for _, ranges := range mergeCellRanges {
			if err := file.MergeCell(sheetName, ranges[0], ranges[1]); err != nil {
				fmt.Println(err)
				return
			}
		}

		headers := []string{"Number", "Student Name"}
		for _, term := range terms {
			headers = append(headers, term.Name)
			headers = append(headers, "%")
		}
		headers = append(headers, "Final Average", "Includes Reinforcement")

		for i, header := range headers {
			cell := string(rune('A'+i)) + "6"
			if header == "Student Name" {
				file.SetCellValue(sheetName, cell, header)
				file.SetCellStyle(sheetName, cell, cell, otherHeaderStyle)
			} else {
				file.SetCellValue(sheetName, cell, header)
				file.SetCellStyle(sheetName, cell, cell, headerStyle)
			}
		}

		_ = file.SetColWidth(subject.Name, "A", "A", 3)
		_ = file.SetColWidth(subject.Name, "B", "B", 37)
		_ = file.SetColWidth(subject.Name, "C", "N", 5)
		_ = file.SetRowHeight(subject.Name, 6, 60)

		for i, student := range students {
			rowNum := i + 7
			file.SetCellValue(sheetName, "A"+strconv.Itoa(rowNum), i+1)
			file.SetCellValue(sheetName, "B"+strconv.Itoa(rowNum), student.Name)

			totalAverage := 0.0
			termCount := 0
			includesReinforcement := false

			for j, term := range terms {
				cellTerm := string(rune('C'+2*j)) + strconv.Itoa(rowNum)
				cellFactor := string(rune('D'+2*j)) + strconv.Itoa(rowNum)
				averageData := getAverageForStudent(averagesData, student.ID, subject.ID, term.Name)
				if averageData != nil {
					file.SetCellValue(sheetName, cellTerm, averageData.Average)
					file.SetCellValue(sheetName, cellFactor, averageData.AveFactor)
					totalAverage += float64(averageData.AveFactor)
					termCount++
					if averageData.Label == "includes_reinforcement" {
						includesReinforcement = true
					}
				} else {
					file.SetCellValue(sheetName, cellTerm, "N/A")
					file.SetCellValue(sheetName, cellFactor, "N/A")
				}
			}

			finalAverageCell := string(rune('C'+2*len(terms))) + strconv.Itoa(rowNum)
			includesReinforcementCell := string(rune('D'+2*len(terms))) + strconv.Itoa(rowNum)
			if termCount > 0 {
				file.SetCellValue(sheetName, finalAverageCell, fmt.Sprintf("%.2f", totalAverage/2))
			} else {
				file.SetCellValue(sheetName, finalAverageCell, "N/A")
			}
			file.SetCellValue(sheetName, includesReinforcementCell, includesReinforcement)
		}
	}

	file.DeleteSheet("Sheet1")

	// Get the current date and time
	currentDatetime := time.Now().Format("2006-01-02_15-04-05")
	// Create the filename in Go format
	filename := fmt.Sprintf("final_averages_report-%s-%s.xlsx", teacherData.TeacherFullName, currentDatetime)

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	//w.Header().Set("Content-Disposition", `attachment; filename="final_averages_report.xlsx"`)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	if err := file.Write(w); err != nil {
		http.Error(w, "Error generating report", http.StatusInternalServerError)
		log.Printf("Error generating report: %v\n", err)
	}
}

func getAverageForStudent(averagesData models.AveragesDataFactor, studentID, subjectID int, termName string) *models.TermAverageFactor {
	for _, studentAverage := range averagesData.Averages {
		if studentAverage.StudentID == studentID && studentAverage.SubjectID == subjectID {
			for _, termAverage := range studentAverage.Averages {
				if termAverage.Term == termName {
					return &termAverage
				}
			}
		}
	}
	return nil
}

///////////////////////////////XLSX TEACHER REPORT///////////////////////

// Handler function to generate the teacher_grades.xlsx report
func GenerateTeacherGradesReport(w http.ResponseWriter, r *http.Request) {
	// Extract the classroom ID, term ID, and teacher ID from the URL path
	vars := mux.Vars(r)
	classroomIDStr := vars["classroomID"]
	termIDStr := vars["termID"]
	teacherIDStr := vars["teacherID"]
	academicPeriodIDStr := vars["academicPeriodID"]

	// Convert classroomID, termID, and teacherID to integers
	classroomID, err := strconv.Atoi(classroomIDStr)
	if err != nil {
		http.Error(w, "Invalid classroom ID", http.StatusBadRequest)
		return
	}

	termID, err := strconv.Atoi(termIDStr)
	if err != nil {
		http.Error(w, "Invalid term ID", http.StatusBadRequest)
		return
	}

	teacherID, err := strconv.Atoi(teacherIDStr)
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	academicPeriodID, err := strconv.Atoi(academicPeriodIDStr)
	if err != nil {
		http.Error(w, "Invalid Academic Period ID", http.StatusBadRequest)
		return
	}

	// Fetch students, subjects, and grades for the classroom and term
	classroom, err := database.GetClassroomByID(classroomID)
	if err != nil {
		http.Error(w, "Error fetching the classroom", http.StatusInternalServerError)
		log.Printf("Error fetching the classroom: %v\n", err)
		return
	}

	students, err := database.GetStudentsByClassroomID(classroomID)
	if err != nil {
		http.Error(w, "Error fetching students", http.StatusInternalServerError)
		log.Printf("Error fetching students: %v\n", err)
		return
	}

	subjects, err := database.GetSubjectsInClassroom(classroomID)
	if err != nil {
		http.Error(w, "Error fetching subjects", http.StatusInternalServerError)
		log.Printf("Error fetching subjects: %v\n", err)
		return
	}

	gradesData, err := database.FetchGradesByClassroomIDAndTermID(classroomID, termID)
	if err != nil {
		http.Error(w, "Error fetching grades", http.StatusInternalServerError)
		log.Printf("Error fetching grades: %v\n", err)
		return
	}

	reinforcementGrades, err := database.GetReinforcementGradeLabelsByClassroomAndTerm(classroomID, termID)
	if err != nil {
		http.Error(w, "Error fetching reinforcement grades", http.StatusInternalServerError)
		log.Printf("Error fetching reinforcement grades: %v\n", err)
		return
	}

	teacherData, err := database.GetTeacherDataByTeacherID(teacherID)
	if err != nil {
		http.Error(w, "Error fetching teacher data", http.StatusInternalServerError)
		log.Printf("Error fetching teacher data: %v\n", err)
		return
	}

	terms, err := database.FetchTermsByAcademicPeriodFromDB(academicPeriodID)
	if err != nil {
		http.Error(w, "Error fetching terms", http.StatusInternalServerError)
		log.Printf("Error fetching terms: %v\n", err)
		return
	}

	var termName string
	validTermID := false
	for _, term := range terms {
		if term.ID == termID {
			validTermID = true
			termName = term.Name
			break
		}
	}

	if !validTermID {
		http.Error(w, "Invalid term ID", http.StatusBadRequest)
		log.Printf("Invalid term ID: %d\n", termID)
		return
	}

	file := excelize.NewFile()

	for _, subject := range subjects {
		sheetName := subject.Name
		file.NewSheet(sheetName)

		gradeLabels, err := database.GetGradeLabelsForSubject(subject.ID, termID)
		if err != nil {
			http.Error(w, "Error fetching grade labels", http.StatusInternalServerError)
			log.Printf("Error fetching grade labels for subject %d: %v\n", subject.ID, err)
			return
		}

		currentTime := time.Now()
		truncatedTime := currentTime.Truncate(time.Second)

		headerSheetStyle, err := file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "center",
			},
			Font: &excelize.Font{
				Bold: true,
			},
		})
		if err != nil {
			fmt.Println(err)
			return
		}

		headerStyle, err := file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Vertical:     "center",
				ShrinkToFit:  true,
				TextRotation: 90,
			},
			Font: &excelize.Font{
				Bold: true,
			},
		})
		if err != nil {
			fmt.Println(err)
			return
		}

		otherHeaderStyle, err := file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Vertical:    "center",
				ShrinkToFit: false,
			},
			Font: &excelize.Font{
				Bold: true,
			},
		})

		if err != nil {
			fmt.Println(err)
			return
		}

		file.SetCellValue(subject.Name, "B1", teacherData.School)
		file.SetCellStyle(sheetName, "B1", "B1", headerSheetStyle)
		file.SetCellValue(subject.Name, "B2", teacherData.Country+" - "+teacherData.City)
		file.SetCellStyle(sheetName, "B2", "B2", headerSheetStyle)
		file.SetCellValue(subject.Name, "B3", "Registro de Calificaciones")
		file.SetCellStyle(sheetName, "B3", "B3", headerSheetStyle)
		file.SetCellValue(subject.Name, "B4", "Teacher: "+teacherData.TeacherFullName)
		file.SetCellValue(subject.Name, "B5", "Classroom: "+classroom.Name)
		file.SetCellValue(subject.Name, "D4", "Subject: "+subject.Name)
		file.SetCellValue(subject.Name, "J4", "School Year: "+teacherData.SchoolYear)
		file.SetCellValue(subject.Name, "D5", "Report Date: "+truncatedTime.Format("2006-01-02 15:04:05"))
		file.SetCellValue(subject.Name, "J5", "School Hours: "+teacherData.SchoolHours)
		file.SetCellValue(subject.Name, "C6", termName)
		file.SetCellStyle(sheetName, "C6", "C6", headerSheetStyle)

		mergeCellRanges := [][]string{{"B1", "O1"}, {"B2", "O2"}, {"B3", "O3"}, {"C6", "O6"}, {"D4", "I4"}, {"J4", "O4"}, {"D5", "I5"}, {"J5", "O5"}}
		for _, ranges := range mergeCellRanges {
			if err := file.MergeCell(sheetName, ranges[0], ranges[1]); err != nil {
				fmt.Println(err)
				return
			}
		}

		headers := []string{"Number", "Student Name"}
		labelIDToName := make(map[int]string)

		for _, label := range gradeLabels {
			headers = append(headers, label.Label)
			labelIDToName[label.ID] = label.Label
		}

		reinforcementLabels := []string{}
		for _, rg := range reinforcementGrades {
			if rg.SubjectID == subject.ID {
				reinforcementLabels = append(reinforcementLabels, rg.Label)
			}
		}
		reinforcementLabels = unique(reinforcementLabels)
		headers = append(headers, reinforcementLabels...)
		headers = append(headers, fmt.Sprintf("%s-average", termName))

		for i, header := range headers {
			cell := string(rune('A'+i)) + "7"
			if header == "Student Name" {
				file.SetCellValue(sheetName, cell, header)
				file.SetCellStyle(sheetName, cell, cell, otherHeaderStyle)
			} else {
				file.SetCellValue(sheetName, cell, header)
				file.SetCellStyle(sheetName, cell, cell, headerStyle)
			}
		}

		_ = file.SetColWidth(subject.Name, "A", "A", 3)
		_ = file.SetColWidth(subject.Name, "B", "B", 37)
		_ = file.SetColWidth(subject.Name, "C", "N", 5)
		_ = file.SetRowHeight(subject.Name, 7, 60)

		for i, student := range students {
			rowNum := i + 8
			file.SetCellValue(sheetName, "A"+strconv.Itoa(rowNum), i+1)
			file.SetCellValue(sheetName, "B"+strconv.Itoa(rowNum), student.Name)

			totalGrades := 0.0
			gradeCount := 0

			for j, label := range gradeLabels {
				cell := string(rune('C'+j)) + strconv.Itoa(rowNum)
				grade := getGradeForStudent(gradesData, student.ID, subject.ID, termName, label.ID)
				if grade != nil {
					file.SetCellValue(sheetName, cell, *grade)
					totalGrades += *grade
					gradeCount++
				} else {
					file.SetCellValue(sheetName, cell, "N/A")
				}
			}

			for k, label := range reinforcementLabels {
				cell := string(rune('C'+len(gradeLabels)+k)) + strconv.Itoa(rowNum)
				grade := getReinforcementGradeForStudent(reinforcementGrades, student.ID, subject.ID, label)
				if grade != nil {
					file.SetCellValue(sheetName, cell, *grade)
					totalGrades += *grade
					gradeCount++
				} else {
					file.SetCellValue(sheetName, cell, "N/A")
				}
			}

			averageCell := string(rune('C'+len(gradeLabels)+len(reinforcementLabels))) + strconv.Itoa(rowNum)
			if gradeCount > 0 {
				file.SetCellValue(sheetName, averageCell, totalGrades/float64(gradeCount))
			} else {
				file.SetCellValue(sheetName, averageCell, "N/A")
			}
		}
	}

	// Get the current date and time
	currentDatetime := time.Now().Format("2006-01-02_15-04-05")
	// Create the filename in Go format
	filename := fmt.Sprintf("teacher-report-%s-%s.xlsx", teacherData.TeacherFullName, currentDatetime)

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	//w.Header().Set("Content-Disposition", "attachment;filename=teacher_grades.xlsx")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	if err := file.Write(w); err != nil {
		http.Error(w, "Error writing file", http.StatusInternalServerError)
		log.Printf("Error writing file: %v\n", err)
		return
	}
}

func getGradeForStudent(gradesData models.GradesData, studentID, subjectID int, termName string, labelID int) *float64 {
	for _, studentGrades := range gradesData.Grades {
		if studentGrades.StudentID == studentID && studentGrades.SubjectID == subjectID {
			for _, termGrades := range studentGrades.Terms {
				if termGrades.Term == termName {
					for _, grade := range termGrades.Grades {
						if grade.LabelID == labelID {
							gradeValue := float64(grade.Grade)
							return &gradeValue
						}
					}
				}
			}
		}
	}
	log.Printf("No grade found for student %d, subject %d, term %s, label ID %d\n", studentID, subjectID, termName, labelID)
	return nil
}

func getReinforcementGradeForStudent(reinforcementGrades []models.ReinforcementGradeLabel, studentID, subjectID int, label string) *float64 {
	for _, grade := range reinforcementGrades {
		if grade.StudentID == studentID && grade.SubjectID == subjectID && grade.Label == label {
			gradeValue := float64(grade.Grade)
			return &gradeValue
		}
	}
	log.Printf("No reinforcement grade found for student %d, subject %d, label %s\n", studentID, subjectID, label)
	return nil
}

func unique(strings []string) []string {
	uniqueStrings := make(map[string]bool)
	for _, str := range strings {
		uniqueStrings[str] = true
	}
	var result []string
	for str := range uniqueStrings {
		result = append(result, str)
	}
	return result
}

/*
///////////////////////////////XLSX TEACHER REPORT WITHOUT REINFORCEMENT///////////////////////

// Handler function to generate the teacher_grades.xlsx report
func GenerateTeacherGradesReport(w http.ResponseWriter, r *http.Request) {
	// Extract the classroom ID, term ID, and teacher ID from the URL path
	vars := mux.Vars(r)
	classroomIDStr := vars["classroomID"]
	termIDStr := vars["termID"]
	teacherIDStr := vars["teacherID"]
	academicPeriodIDStr := vars["academicPeriodID"]

	// Convert classroomID, termID, and teacherID to integers
	classroomID, err := strconv.Atoi(classroomIDStr)
	if err != nil {
		http.Error(w, "Invalid classroom ID", http.StatusBadRequest)
		return
	}

	termID, err := strconv.Atoi(termIDStr)
	if err != nil {
		http.Error(w, "Invalid term ID", http.StatusBadRequest)
		return
	}

	teacherID, err := strconv.Atoi(teacherIDStr)
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	academicPeriodID, err := strconv.Atoi(academicPeriodIDStr)
	if err != nil {
		http.Error(w, "Invalid Academic Period ID", http.StatusBadRequest)
		return
	}

	// Fetch students, subjects, and grades for the classroom and term
	classroom, err := database.GetClassroomByID(classroomID)
	if err != nil {
		http.Error(w, "Error fetching the classroom", http.StatusInternalServerError)
		log.Printf("Error fetching the classroom: %v\n", err)
		return
	}

	// Fetch students, subjects, and grades for the classroom and term
	students, err := database.GetStudentsByClassroomID(classroomID)
	if err != nil {
		http.Error(w, "Error fetching students", http.StatusInternalServerError)
		log.Printf("Error fetching students: %v\n", err)
		return
	}

	subjects, err := database.GetSubjectsInClassroom(classroomID)
	if err != nil {
		http.Error(w, "Error fetching subjects", http.StatusInternalServerError)
		log.Printf("Error fetching subjects: %v\n", err)
		return
	}

	gradesData, err := database.FetchGradesByClassroomIDAndTermID(classroomID, termID)
	if err != nil {
		http.Error(w, "Error fetching grades", http.StatusInternalServerError)
		log.Printf("Error fetching grades: %v\n", err)
		return
	}

	teacherData, err := database.GetTeacherDataByTeacherID(teacherID)
	if err != nil {
		http.Error(w, "Error fetching teacher data", http.StatusInternalServerError)
		log.Printf("Error fetching teacher data: %v\n", err)
		return
	}

	// Fetch terms by teacher ID to validate term
	terms, err := database.FetchTermsByAcademicPeriodFromDB(academicPeriodID)
	if err != nil {
		http.Error(w, "Error fetching terms", http.StatusInternalServerError)
		log.Printf("Error fetching terms: %v\n", err)
		return
	}

	// Check if the term ID is valid and get the term name
	var termName string
	validTermID := false
	for _, term := range terms {
		if term.ID == termID {
			validTermID = true
			termName = term.Name
			break
		}
	}

	if !validTermID {
		http.Error(w, "Invalid term ID", http.StatusBadRequest)
		log.Printf("Invalid term ID: %d\n", termID)
		return
	}

	// Generate the Excel file
	file := excelize.NewFile()

	// Generate a separate sheet for each subject
	for _, subject := range subjects {
		sheetName := subject.Name
		file.NewSheet(sheetName)

		// Fetch grade labels for the current subject and term
		gradeLabels, err := database.GetGradeLabelsForSubject(subject.ID, termID)
		if err != nil {
			http.Error(w, "Error fetching grade labels", http.StatusInternalServerError)
			log.Printf("Error fetching grade labels for subject %d: %v\n", subject.ID, err)
			return
		}
		// Get the current date and time
		currentTime := time.Now()
		// Truncate the time to seconds
		truncatedTime := currentTime.Truncate(time.Second)

		// Apply styles to the header sheet
		headerSheetStyle, err := file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "center",
			},
			Font: &excelize.Font{
				Bold: true,
			},
		})
		if err != nil {
			fmt.Println(err)
			return
		}

		// Apply styles to the header row
		headerStyle, err := file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Vertical:     "center",
				ShrinkToFit:  true,
				TextRotation: 90, // Set to 90 degrees for vertical text
			},
			Font: &excelize.Font{
				Bold: true,
			},
		})
		if err != nil {
			fmt.Println(err)
			return
		}

		// Apply styles to other headers (excluding "Student Name")
		otherHeaderStyle, err := file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Vertical:    "center",
				ShrinkToFit: false,
			},
			Font: &excelize.Font{
				Bold: true,
			},
		})

		if err != nil {
			fmt.Println(err)
			return
		}
		// Set the header of the sheet
		file.SetCellValue(subject.Name, "B1", teacherData.School)
		file.SetCellStyle(sheetName, "B1", "B1", headerSheetStyle)
		file.SetCellValue(subject.Name, "B2", teacherData.Country+" - "+teacherData.City)
		file.SetCellStyle(sheetName, "B2", "B2", headerSheetStyle)
		file.SetCellValue(subject.Name, "B3", "Teacher: "+teacherData.TeacherFullName)
		file.SetCellValue(subject.Name, "B4", "Classroom: "+classroom.Name)
		file.SetCellValue(subject.Name, "D3", "Subject: "+subject.Name)
		file.SetCellValue(subject.Name, "J3", "School Year: "+teacherData.SchoolYear)
		file.SetCellValue(subject.Name, "D4", "Report Date: "+truncatedTime.Format("2006-01-02 15:04:05"))
		file.SetCellValue(subject.Name, "J4", "School Hours: "+teacherData.SchoolHours)
		file.SetCellValue(subject.Name, "C5", termName)
		file.SetCellStyle(sheetName, "C5", "C5", headerSheetStyle)

		// Merge cells
		mergeCellRanges := [][]string{{"B1", "O1"}, {"B2", "O2"}, {"C5", "O5"}, {"D3", "I3"}, {"J3", "O3"}, {"D4", "I4"}, {"J4", "O4"}}
		for _, ranges := range mergeCellRanges {
			if err := file.MergeCell(sheetName, ranges[0], ranges[1]); err != nil {
				fmt.Println(err)
				return
			}
		}

		// Set the header labels row
		headers := []string{"Number", "Student Name"}
		labelIDToName := make(map[int]string)

		for _, label := range gradeLabels {
			headers = append(headers, label.Label) // Add label names from gradeLabels
			labelIDToName[label.ID] = label.Label
		}
		headers = append(headers, fmt.Sprintf("%s-average", termName))

		for i, header := range headers {
			cell := string(rune('A'+i)) + "6"
			if header == "Student Name" {
				file.SetCellValue(sheetName, cell, header)
				file.SetCellStyle(sheetName, cell, cell, otherHeaderStyle)
			} else {
				file.SetCellValue(sheetName, cell, header)
				file.SetCellStyle(sheetName, cell, cell, headerStyle)
			}
		}

		_ = file.SetColWidth(subject.Name, "A", "A", 3)
		_ = file.SetColWidth(subject.Name, "B", "B", 37)
		_ = file.SetColWidth(subject.Name, "C", "N", 5)
		_ = file.SetRowHeight(subject.Name, 6, 60)

		// Fill the student grades
		for i, student := range students {
			rowNum := i + 7
			file.SetCellValue(sheetName, "A"+strconv.Itoa(rowNum), i+1)
			file.SetCellValue(sheetName, "B"+strconv.Itoa(rowNum), student.Name)

			totalGrades := 0.0
			gradeCount := 0

			for j, label := range gradeLabels {
				cell := string(rune('C'+j)) + strconv.Itoa(rowNum)
				grade := getGradeForStudent(gradesData, student.ID, subject.ID, termName, label.ID)
				if grade != nil {
					file.SetCellValue(sheetName, cell, *grade)
					totalGrades += *grade
					gradeCount++
				} else {
					file.SetCellValue(sheetName, cell, "N/A")
				}
			}

			averageCell := string(rune('C'+len(gradeLabels))) + strconv.Itoa(rowNum)
			if gradeCount > 0 {
				file.SetCellValue(sheetName, averageCell, totalGrades/float64(gradeCount))
			} else {
				file.SetCellValue(sheetName, averageCell, "N/A")
			}
		}
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment;filename=teacher_grades.xlsx")

	// Write the file to the response
	if err := file.Write(w); err != nil {
		http.Error(w, "Error writing file", http.StatusInternalServerError)
		log.Printf("Error writing file: %v\n", err)
		return
	}
}

// Helper function to get the grade for a student in a specific subject and term
func getGradeForStudent(gradesData models.GradesData, studentID, subjectID int, termName string, labelID int) *float64 {
	for _, studentGrades := range gradesData.Grades {
		if studentGrades.StudentID == studentID && studentGrades.SubjectID == subjectID {
			for _, termGrades := range studentGrades.Terms {
				if termGrades.Term == termName {
					for _, grade := range termGrades.Grades {
						if grade.LabelID == labelID {
							gradeValue := float64(grade.Grade) // Convert float32 to float64
							return &gradeValue
						}
					}
				}
			}
		}
	}
	log.Printf("No grade found for student %d, subject %d, term %s, label ID %d\n", studentID, subjectID, termName, labelID)
	return nil
}
*/

/*
err = file.MergeCell(sheetName, "A1", "N1")
			if err != nil {
				fmt.Println(err)
				return
			}

// Set value of a cell.
		file.SetCellValue("English", "H7", "Hello, World!")

		// Set vertical alignment to center and shrink text to fit the cell.
		style, err := file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Vertical:     "center", // Set to "top" or "bottom" as needed
				ShrinkToFit:  true,
				WrapText:     true,
				TextRotation: 90, // Set to 90 degrees for vertical text
			},
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		file.SetCellStyle("English", "H7", "H7", style)
*/
