//  static/classroom-grades-upload.js

document.addEventListener("DOMContentLoaded", async () => {
    const urlParams = new URLSearchParams(window.location.search);
    const classroomID = urlParams.get('classroomID');
    const classroomName = urlParams.get('classroomName');
    
    document.getElementById("classroom-title").textContent = `Grades Grid of ${classroomName}`;

    try {
        const [studentsResponse, subjectsResponse, gradesResponse] = await Promise.all([
            fetch(`/classrooms/${classroomID}/students`),
            fetch(`/classrooms/${classroomID}/subjects`),
            fetch(`/classrooms/${classroomID}/grades/get`)
        ]);

        const students = await studentsResponse.json();
        const subjects = await subjectsResponse.json();
        const gradesData = await gradesResponse.json();

        const gradesGridContainer = document.getElementById("grades-grid-container");

        subjects.forEach(subject => {
            const subjectContainer = document.createElement('div');
            subjectContainer.classList.add('subject-container');
            
            const subjectTitle = document.createElement('h2');
            subjectTitle.textContent = subject.name;
            subjectContainer.appendChild(subjectTitle);

            const table = document.createElement('table');
            table.id = `grades-grid-${subject.id}`;
            table.classList.add('grades-grid');
            subjectContainer.appendChild(table);

            const subjectGrades = gradesData.grades.filter(grade => grade.subjectID === subject.id);  // Add grades data for the current subject
            generateGradesGrid(table, students, subject.grade_labels, subjectGrades);
            
            gradesGridContainer.appendChild(subjectContainer);
            //generateGradesGrid(table, students, subject.grade_labels);
        });

         // Add event listener to the upload button
         document.getElementById("upload-grades-btn").addEventListener("click", () => {
            uploadGrades(classroomID, students, subjects);
        });
    } catch (error) {
        console.error('Error fetching students or subjects:', error);
    }
});

function generateGradesGrid(gridElement, students, gradeLabels, gradesData) {
    const headerRow = gridElement.insertRow();

    const numberHeaderCell = headerRow.insertCell();
    numberHeaderCell.textContent = 'Number';

    const nameHeaderCell = headerRow.insertCell();
    nameHeaderCell.textContent = 'Student Name';

    gradeLabels.forEach(label => {
        const headerCell = headerRow.insertCell();
        headerCell.textContent = label;
    });

    students.forEach((student, index) => {
        const row = gridElement.insertRow();

        const numberCell = row.insertCell();
        numberCell.textContent = index + 1;

        const nameCell = row.insertCell();
        nameCell.textContent = student.name;
        gradeLabels.forEach(label => {
            const gradeCell = row.insertCell();
            gradeCell.contentEditable = true;

            const studentGrade = gradesData.find(grade => grade.studentID === student.id);
            if (studentGrade) {
                const gradeItem = studentGrade.grades.find(grade => grade.label === label);
                if (gradeItem) {
                    gradeCell.textContent = gradeItem.grade;
                }
            }
        });
    });
}

function uploadGrades(classroomID, students, subjects) {
    const gradesData = [];

    subjects.forEach(subject => {
        const gradesGrid = document.getElementById(`grades-grid-${subject.id}`);

        for (let i = 1; i < gradesGrid.rows.length; i++) {
            const row = gradesGrid.rows[i];
            const studentGrades = {
                studentID: students[i - 1].id,
                subjectID: subject.id,
                grades: []
            };

            for (let j = 2; j < row.cells.length; j++) {
                studentGrades.grades.push({
                    label: subject.grade_labels[j - 2],
                    grade: row.cells[j].textContent.trim()
                });
            }

            gradesData.push(studentGrades);
        }
    });

    console.log("Grades Data to be uploaded:", JSON.stringify({ grades: gradesData }));

    fetch(`/classrooms/${classroomID}/grades`, { // Updated URL with classroomID
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ grades: gradesData })
    })
    .then(response => {
        console.log("Response from server:", response);
        if (!response.ok) {
            throw new Error('Failed to upload grades');
        }
        return response.json();
    })
    .then(data => {
        console.log("Data received from server:", data);
        alert('Grades uploaded successfully');
    })
    .catch(error => {
        console.error('Error uploading grades:', error);
        alert('Failed to upload grades');
    });
}
