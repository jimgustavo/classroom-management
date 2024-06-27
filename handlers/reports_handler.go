// handlers/reports_handlers.go
package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

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

func loadTemplate(pdf *gofpdf.Fpdf, teacherData *models.TeacherData, classroom *models.Classroom, student *models.Student) error {

	//In case is needed to use special characters such as: ñ, é, ó and so forth.
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	// Add logo image
	logoPath := "ue12f_logo.jpeg"
	pdf.Image(logoPath, 10, 5, 20, 0, false, "", 0, "ue12f_logo")

	// Your template content
	schoolName := teacherData.School
	studentName := student.Name
	classroomName := classroom.Name
	//schoolYear := teacherData.SchoolYear
	//countryCityDate := fmt.Sprintf("%s, %s %s", teacherData.Country, teacherData.City, time.Now().Format("02 January 2006"))

	// Adding content to the document
	// Add title
	pdf.CellFormat(180, 10, schoolName, "0", 0, "C", false, 0, "")
	pdf.Ln(10)
	pdf.Cell(40, 30, tr("ACTA DE COMPROMISO POR BAJO RENDIMIENTO ACADéMICO"))
	pdf.Ln(20)

	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(0, 10, fmt.Sprintf(tr("En la ciudad de %s, a los %s días del mes de %s del %s, comparecen ante el/la rector(a)/vicerrector(a) de la %s, el/la Sr./Sra. ____________________________, en calidad de representante legal del estudiante %s del curso %s, para suscribir la presente acta de compromiso."),
		teacherData.City, time.Now().Format("02"), time.Now().Format("January"), time.Now().Format("2006"), teacherData.School, studentName, classroomName), "", "", false)
	pdf.Ln(10)
	pdf.MultiCell(0, 10, tr("Considerando:"), "", "", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 10, tr("1. Que, conforme al Artículo 26 de la Constitución de la República del Ecuador, \"la educación es un derecho de las personas a lo largo de su vida y un deber ineludible e inexcusable del Estado. Constituye un área prioritaria de la política pública y de la inversión estatal, garantía de la igualdad e inclusión social y condición indispensable para el buen vivir. Las personas, las familias y la sociedad tienen el derecho y la responsabilidad de participar en el proceso educativo.\""), "", "", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 10, tr("2. Que, de acuerdo con el Artículo 8 de la Ley Orgánica de Educación Intercultural (LOEI), las y los estudiantes tienen obligaciones y responsabilidades, tales como cumplir con las actividades académico-formativas, participar en evaluaciones, procurar la excelencia educativa y mostrar integridad y honestidad académica."), "", "", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 10, tr("3. Que, según el Artículo 13 de la LOEI, las madres, padres y/o representantes de los estudiantes deben involucrarse activamente en los procesos educativos de sus representados y atender los llamados y requerimientos de los profesores y autoridades de los planteles, y apoyar y motivar a sus representados especialmente cuando existan dificultades en el proceso de aprendizaje."), "", "", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 10, tr("4. Que, el Artículo 32 del Reglamento a la Ley Orgánica de Educación Intercultural (RLOEI) establece que, si la evaluación continua determinare bajos resultados en los procesos de aprendizaje, se deberá diseñar e implementar de inmediato procesos de refuerzo pedagógico."), "", "", false)
	pdf.Ln(10)

	pdf.MultiCell(0, 10, tr("Compromisos:"), "", "", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 10, fmt.Sprintf(tr("El/la Sr./Sra. ________________________________, en su calidad de representante legal del estudiante %s, se compromete a:"), studentName), "", "", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 10, fmt.Sprintf(tr("1. Garantizar la asistencia regular del estudiante %s a todas las clases y actividades educativas programadas."), studentName), "", "", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 10, tr("2. Colaborar activamente con el personal docente en el diseño e implementación de un plan de refuerzo pedagógico, que incluirá clases de refuerzo, tutorías individuales y un cronograma de estudios a seguir en casa."), "", "", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 10, tr("3. Proveer un ambiente de aprendizaje adecuado en el hogar, dedicando espacios específicos para las tareas y estudios del estudiante."), "", "", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 10, tr("4. Motivar y apoyar al estudiante en el cumplimiento de sus obligaciones académicas, fomentando la excelencia educativa y la integridad académica."), "", "", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 10, tr("5. Participar en reuniones de seguimiento con los docentes y autoridades del plantel para evaluar el progreso académico del estudiante."), "", "", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 10, tr("6. Reconocer y valorar los esfuerzos y avances del estudiante, así como los méritos y la excelencia del personal docente."), "", "", false)
	pdf.Ln(10)

	// Calificaciones actuales del estudiante
	pdf.MultiCell(0, 10, tr("Calificaciones actuales del estudiante:"), "", "", false)
	pdf.Ln(5)

	pdf.MultiCell(0, 10, tr("Firman en conformidad:"), "", "", false)
	pdf.Ln(5)

	signatureLines := []string{
		"Representante Legal",
		"Estudiante",
		"Rector(a)/Vicerrector(a)",
		"Docente Tutor/a",
	}
	for _, line := range signatureLines {
		pdf.MultiCell(0, 10, tr(line), "", "C", false)
		pdf.Ln(5)
	}
	pdf.Ln(10)

	pdf.MultiCell(0, 10, fmt.Sprintf("CI: %s", strconv.Itoa(teacherData.TeacherID)), "", "C", false)
	pdf.MultiCell(0, 10, fmt.Sprintf("Correo: %s", teacherData.InstitutionalEmail), "", "C", false)
	pdf.MultiCell(0, 10, fmt.Sprintf("Telefono: %s", teacherData.Phone), "", "C", false)
	pdf.Ln(10)

	return nil
}

//////////////////////////////XLSX TERM REPORT///////////////////////////

func GenerateAveragesExcelReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	classroomIDStr := vars["classroomID"]

	classroomID, err := strconv.Atoi(classroomIDStr)
	if err != nil {
		http.Error(w, "Invalid classroom ID", http.StatusBadRequest)
		return
	}

	// Fetch students, subjects, and averages
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

	// Dummy termFactors. Adjust as needed or fetch dynamically if required.
	termFactors := []models.TermFactor{
		{Term: "bimestre1", Factor: 0.8},
		{Term: "bimestre2", Factor: 0.2},
	}

	averagesData, err := database.FetchAveragesWithFactorsByClassroomID(classroomID, termFactors)
	if err != nil {
		http.Error(w, "Error fetching averages", http.StatusInternalServerError)
		log.Printf("Error fetching averages: %v\n", err)
		return
	}

	// Generate the Excel file
	file := excelize.NewFile()

	// Generate a separate sheet for each subject
	for _, subject := range subjects {
		sheetName := subject.Name
		file.NewSheet(sheetName)

		// Set the header row
		headers := []string{"Number", "Student Name"}
		terms := []string{"bimestre1", "bimestre2"}

		for _, term := range terms {
			headers = append(headers, term)
			headers = append(headers, fmt.Sprintf("%s-%s", term, "Average-%"))
		}
		headers = append(headers, "Final Average")

		for i, header := range headers {
			cell := string(rune('A'+i)) + "1"
			file.SetCellValue(sheetName, cell, header)
		}

		// Fill the student grades
		for i, student := range students {
			rowNum := i + 2
			file.SetCellValue(sheetName, "A"+strconv.Itoa(rowNum), i+1)
			file.SetCellValue(sheetName, "B"+strconv.Itoa(rowNum), student.Name)

			totalAverage := 0.0
			termCount := 0

			for j, term := range terms {
				cell := string(rune('C'+2*j)) + strconv.Itoa(rowNum)
				factorCell := string(rune('C'+2*j+1)) + strconv.Itoa(rowNum)

				studentAverage := getAverageForStudent(averagesData, student.ID, subject.ID, term)
				if studentAverage != nil {
					file.SetCellValue(sheetName, cell, studentAverage.Average)
					file.SetCellValue(sheetName, factorCell, studentAverage.AveFactor)
					totalAverage += float64(studentAverage.AveFactor) // Convert to float64
					termCount++
				} else {
					file.SetCellValue(sheetName, cell, "N/A")
					file.SetCellValue(sheetName, factorCell, "N/A")
				}
			}

			averageCell := string(rune('C'+2*len(terms))) + strconv.Itoa(rowNum)
			if termCount > 0 {
				//file.SetCellValue(sheetName, averageCell, totalAverage/float64(termCount))
				file.SetCellValue(sheetName, averageCell, totalAverage)
			} else {
				file.SetCellValue(sheetName, averageCell, "N/A")
			}
		}
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment;filename=averages.xlsx")

	// Write the file to the response
	if err := file.Write(w); err != nil {
		http.Error(w, "Error writing file", http.StatusInternalServerError)
		log.Printf("Error writing file: %v\n", err)
		return
	}
}

// Helper function to get the average for a student in a specific subject and term
func getAverageForStudent(averagesData models.AveragesDataFactor, studentID, subjectID int, term string) *models.TermAverageFactor {
	for _, studentAverages := range averagesData.Averages {
		if studentAverages.StudentID == studentID && studentAverages.SubjectID == subjectID {
			for _, termAverage := range studentAverages.Averages {
				if termAverage.Term == term {
					return &termAverage
				}
			}
		}
	}
	log.Printf("No average found for student %d, subject %d, term %s\n", studentID, subjectID, term)
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
	terms, err := database.GetTermsByTeacherID(teacherID)
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
