//  static/classroom-grades-display.js

document.addEventListener("DOMContentLoaded", async () => {
    const urlParams = new URLSearchParams(window.location.search);
    const classroomID = urlParams.get('classroomID');
    const termID = urlParams.get('termID');
    const term = urlParams.get('term');
    const classroomName = urlParams.get('classroomName');
    
    document.getElementById("classroom-title").textContent = `Grades Grid of ${classroomName} - ${term}`;

    try {
        const [studentsResponse, subjectsResponse, gradesResponse] = await Promise.all([
            fetch(`/classrooms/${classroomID}/students`),
            fetch(`/classrooms/${classroomID}/subjects`),
            fetch(`/classrooms/${classroomID}/terms/${termID}/grades`)
        ]);

        const students = await studentsResponse.json();
        const subjects = await subjectsResponse.json();
        const gradesData = await gradesResponse.json();

        console.log("Fetched Students:", students);
        console.log("Fetched Subjects:", subjects);
        console.log("Fetched Grades Data:", gradesData);
       
        const gradesGridContainer = document.getElementById("grades-grid-container");

        subjects.forEach(subject => {
            const filteredGradeLabels = subject.grade_labels.filter(label => label.term_id == termID); // Filter by termID
            if (filteredGradeLabels.length > 0) {
                const subjectContainer = document.createElement('div');
                subjectContainer.classList.add('subject-container');

                const subjectTitle = document.createElement('h2');
                subjectTitle.textContent = subject.name;
                subjectContainer.appendChild(subjectTitle);

                const table = document.createElement('table');
                table.id = `grades-grid-${subject.id}`;
                table.classList.add('grades-grid');
                subjectContainer.appendChild(table);

                generateGradesGrid(table, students, filteredGradeLabels, gradesData, subject.id, term);

                gradesGridContainer.appendChild(subjectContainer);
            }
        });
    } catch (error) {
        console.error('Error fetching students, subjects, or grades:', error);
    }
});

function generateGradesGrid(gridElement, students, gradeLabels, gradesData, subjectID, term) {
    const headerRow = gridElement.insertRow();

    const numberHeaderCell = headerRow.insertCell();
    numberHeaderCell.textContent = 'Number';

    const nameHeaderCell = headerRow.insertCell();
    nameHeaderCell.textContent = 'Student Name';

    gradeLabels.forEach(label => {
        const headerCell = headerRow.insertCell();
        headerCell.textContent = label.label; // Display label name
    });

    students.forEach((student, index) => {
        const row = gridElement.insertRow();

        const numberCell = row.insertCell();
        numberCell.textContent = index + 1;

        const nameCell = row.insertCell();
        nameCell.textContent = student.name;

        gradeLabels.forEach(label => {
            const gradeCell = row.insertCell();
            gradeCell.contentEditable = false;
            gradeCell.dataset.labelId = label.id; // Store label ID in data attribute

            // Find the corresponding grade for this student, subject, and label
            const studentGrades = gradesData.grades.find(g => g.student_id === student.id && g.subject_id === subjectID);
            if (studentGrades) {
                const termGrades = studentGrades.terms.find(t => t.term === term);
                if (termGrades) {
                    const gradeEntry = termGrades.grades.find(g => g.label_id === label.id);
                    if (gradeEntry) {
                        gradeCell.textContent = gradeEntry.grade;
                    }
                }
            }
        });
    });
}


/*
function generateGradesGrid(gridElement, students, gradeLabels) {
    const headerRow = gridElement.insertRow();

    const numberHeaderCell = headerRow.insertCell();
    numberHeaderCell.textContent = 'Number';

    const nameHeaderCell = headerRow.insertCell();
    nameHeaderCell.textContent = 'Student Name';

    gradeLabels.forEach(label => {
        const headerCell = headerRow.insertCell();
        headerCell.textContent = label.label; // Display label name
    });

    // Log the grade labels
    console.log('Grade Labels:', gradeLabels);

    students.forEach((student, index) => {
        const row = gridElement.insertRow();

        const numberCell = row.insertCell();
        numberCell.textContent = index + 1;

        const nameCell = row.insertCell();
        nameCell.textContent = student.name;

        gradeLabels.forEach(label => {
            const gradeCell = row.insertCell();
            gradeCell.contentEditable = false; // Make cell editable

            // Find the student's grade for the current label
            const studentGrade = subjectGrades.find(grade => grade.student_id === student.id);
            
            // Log the student grade
            console.log(`Student ID: ${student.id}, Student Grade:`, studentGrade);

            if (studentGrade) {
                const gradeItem = studentGrade.grades.find(g => g.label_id === label.id);

                // Log the grade item
                console.log(`Label ID: ${label.id}, Grade Item:`, gradeItem);

                if (gradeItem) {
                    gradeCell.textContent = gradeItem.grade;
                }
            }
        });
    });
}

*/

 /*
        // Check if gradesData is empty
        if (!gradesData || gradesData.length === 0) {
            console.log("No grades data found.");
            gradesData = { grades: [] }; // Set gradesData to an empty object
        }
           
        // Log label_id and grade
        gradesData.grades.forEach(student => {
            student.terms.forEach(term => {
                 console.log(`Term: ${term.term}`);
                 term.grades.forEach(grade => {
                     console.log(`  label_id: ${grade.label_id}, grade: ${grade.grade}`);
                });
            });
        });
        */


         /*
        const subjectGrades = gradesData.grades
                .filter(grade => grade.subject_id === subject.id)
                .flatMap(grade => grade.terms)
                .filter(t => t.term === term); // Filter grades by term
         */

/*
            // Get the grades for the current student
            const studentTermGrades = gradesData.find(grade => grade.student_id === student.id);

            if (studentTermGrades) {
                const gradeItem = studentTermGrades.grades.find(grade => grade.label_id === label.id); // Adjusted to use label_id
                if (gradeItem) {
                    gradeCell.textContent = gradeItem.grade;
                }
            }
             */
