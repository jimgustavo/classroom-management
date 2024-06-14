//  static/classroom-grades-upload.js

document.addEventListener("DOMContentLoaded", async () => {
    const urlParams = new URLSearchParams(window.location.search);
    const classroomID = urlParams.get('classroomID');
    const term = urlParams.get('term');
    const termID = urlParams.get('termID');
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
        //const gradesData = await gradesResponse.json();

        //console.log("Fetched Students:", students);
        //console.log("Fetched Subjects:", subjects);
        //console.log("Fetched Grades Data:", gradesData);

        const gradesGridContainer = document.getElementById("grades-grid-container");

        subjects.forEach(subject => {
            const filteredGradeLabels = subject.grade_labels.filter(label => label.term_id == termID); // Filter by termID
            console.log(`subjectID: ${subject.id}, filteredLabels:`, filteredGradeLabels);
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

                generateGradesGrid(table, students, filteredGradeLabels);

                gradesGridContainer.appendChild(subjectContainer);
            }
        });

        document.getElementById("upload-grades-btn").addEventListener("click", () => {
            uploadGrades(classroomID, students, subjects, term);
        });

    } catch (error) {
        console.error('Error fetching students or subjects:', error);
    }
});

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

    students.forEach((student, index) => {
        const row = gridElement.insertRow();

        const numberCell = row.insertCell();
        numberCell.textContent = index + 1;

        const nameCell = row.insertCell();
        nameCell.textContent = student.name;

        gradeLabels.forEach(label => {
            const gradeCell = row.insertCell();
            gradeCell.contentEditable = true;
            gradeCell.dataset.labelId = label.id; // Store label ID in data attribute
        });
    });
}

function uploadGrades(classroomID, students, subjects, term) {
    const gradesData = [];

    subjects.forEach(subject => {
        const gradesGrid = document.getElementById(`grades-grid-${subject.id}`);

        for (let i = 1; i < gradesGrid.rows.length; i++) {
            const row = gradesGrid.rows[i];
            const studentGrades = {
                student_id: students[i - 1].id,
                subject_id: subject.id,
                terms: [
                    {
                        term: term,
                        grades: []
                    }
                ]
            };

            for (let j = 2; j < row.cells.length; j++) {
                studentGrades.terms[0].grades.push({
                    label_id: parseInt(row.cells[j].dataset.labelId),
                    grade: parseFloat(row.cells[j].textContent.trim()) || 0
                });
            }

            gradesData.push(studentGrades);
        }
    });

    console.log("Grades Data to be uploaded:", JSON.stringify({ grades: gradesData }));

    fetch(`/api/classrooms/${classroomID}/grades`, {
        method: 'POST',
        headers: {
            "Content-Type": "application/json",
            'Authorization': `Bearer ${localStorage.getItem("token")}`
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
        console.log("uploaded grades:", body);
    })
    .catch(error => {
        console.error('Error uploading grades:', error);
        alert('Failed to upload grades');
    });
}


 /*
            for (let j = 2; j < row.cells.length - 1; j++) { // Exclude the last cell (term-average)
                studentGrades.terms[0].grades.push({
                    label_id: parseInt(row.cells[j].dataset.labelId),
                    grade: parseFloat(row.cells[j].textContent.trim()) || 0
                });
            }
*/