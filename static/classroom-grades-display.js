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
            fetch(`/api/classrooms/${classroomID}/students`, {
                method: 'GET', // Add the GET method
                headers: {
                    "Content-Type": "application/json",
                    'Authorization': `Bearer ${localStorage.getItem("token")}` // Add your authorization token here
                }
            }),
            fetch(`/api/classrooms/${classroomID}/subjects`, {
                method: 'GET', // Add the GET method
                headers: {
                    "Content-Type": "application/json",
                    'Authorization': `Bearer ${localStorage.getItem("token")}` // Add your authorization token here
                }
            }),
            fetch(`/api/classrooms/${classroomID}/grades/get?term=${encodeURIComponent(term)}`, {
                method: 'GET', // Add the GET method
                headers: {
                    "Content-Type": "application/json",
                    'Authorization': `Bearer ${localStorage.getItem("token")}` // Add your authorization token here
                }
            })
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

    // Add term-average header
    const averageHeaderCell = headerRow.insertCell();
    averageHeaderCell.textContent = `${term}-average`;
    
    students.forEach((student, index) => {
        const row = gridElement.insertRow();

        const numberCell = row.insertCell();
        numberCell.textContent = index + 1;

        const nameCell = row.insertCell();
        nameCell.textContent = student.name;

        let totalGrades = 0;
        let gradeCount = 0;

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
                        totalGrades += parseFloat(gradeEntry.grade);
                        gradeCount++;
                    }
                }
            }
        });

            // Calculate and add term average
            const averageCell = row.insertCell();
            const average = gradeCount > 0 ? (totalGrades / gradeCount).toFixed(2) : '0.00';
            averageCell.textContent = average;
    });
}