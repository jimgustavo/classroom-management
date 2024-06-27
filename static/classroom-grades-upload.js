//  static/classroom-grades-upload.js

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
                method: 'GET',
                headers: {
                    "Content-Type": "application/json",
                    'Authorization': `Bearer ${localStorage.getItem("token")}`
                }
            }),
            fetch(`/api/classrooms/${classroomID}/subjects`, {
                method: 'GET',
                headers: {
                    "Content-Type": "application/json",
                    'Authorization': `Bearer ${localStorage.getItem("token")}`
                }
            }),
            fetch(`/api/classrooms/${classroomID}/terms/${termID}/grades`, {
                method: 'GET',
                headers: {
                    "Content-Type": "application/json",
                    'Authorization': `Bearer ${localStorage.getItem("token")}`
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

                //generateGradesGrid(table, students, filteredGradeLabels, gradesData, subject.id, term);
                generateGradesGrid(table, students, filteredGradeLabels, gradesData.grades || [], subject.id, term);
                
                gradesGridContainer.appendChild(subjectContainer);
            }
        });

        document.getElementById("upload-grades-btn").addEventListener("click", () => {
            uploadGrades(classroomID, students, subjects, term);
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
            gradeCell.contentEditable = true; // Allow editing
            gradeCell.dataset.labelId = label.id; // Store label ID in data attribute
            gradeCell.dataset.studentId = student.id; // Store student ID in data attribute
            gradeCell.dataset.subjectId = subjectID; // Store subject ID in data attribute
            gradeCell.dataset.term = term; // Store term in data attribute

            // Find the corresponding grade for this student, subject, and label
          //const studentGrades = gradesData.grades.find(g => g.student_id === student.id && g.subject_id === subjectID);
            const studentGrades = gradesData.find(g => g.student_id === student.id && g.subject_id === subjectID);
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

async function uploadGrades(classroomID, students, subjects, term) {
    const token = localStorage.getItem("token");
    const gradesData = { grades: [] };

    // Create a mapping to group grades by student and subject
    const gradesMap = {};

    students.forEach(student => {
        gradesMap[student.id] = gradesMap[student.id] || {};
        subjects.forEach(subject => {
            gradesMap[student.id][subject.id] = gradesMap[student.id][subject.id] || { term: term, grades: [] };
        });
    });

    subjects.forEach(subject => {
        const table = document.getElementById(`grades-grid-${subject.id}`);
        if (!table) return;

        const rows = table.rows;

        for (let i = 1; i < rows.length; i++) { // Skip the header row
            const cells = rows[i].cells;
            const studentID = students[i - 1].id; // Adjust index since we skip the header row

            for (let j = 2; j < cells.length - 1; j++) { // Skip the number and name columns and the average column
                const cell = cells[j];
                const grade = parseFloat(cell.textContent);
                const labelID = parseInt(cell.dataset.labelId);

                if (!isNaN(grade)) { // Only add valid grades
                    gradesMap[studentID][subject.id].grades.push({
                        label_id: labelID,
                        grade: grade
                    });
                }
            }
        }
    });

    // Convert the gradesMap to the expected structure
    for (const studentID in gradesMap) {
        for (const subjectID in gradesMap[studentID]) {
            gradesData.grades.push({
                student_id: parseInt(studentID),
                subject_id: parseInt(subjectID),
                terms: [gradesMap[studentID][subjectID]]
            });
        }
    }

    console.log("Updated Grades Payload:", JSON.stringify(gradesData));

    try {
        const response = await fetch(`/api/classrooms/${classroomID}/grades`, {
            method: 'POST',
            headers: {
                "Content-Type": "application/json",
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(gradesData)
        });

        if (response.ok) {
            alert("Grades updated successfully!");
        } else {
            console.error('Error updating grades:', response.statusText);
            alert("Error updating grades. Please try again.");
        }
    } catch (error) {
        console.error('Error updating grades:', error);
        alert("Error updating grades. Please try again.");
    }
}

/*
document.addEventListener("DOMContentLoaded", async () => {
    const urlParams = new URLSearchParams(window.location.search);
    const classroomID = urlParams.get('classroomID');
    const term = urlParams.get('term');
    const termID = urlParams.get('termID');
    const classroomName = urlParams.get('classroomName');

    document.getElementById("classroom-title").textContent = `Grades Grid of ${classroomName} - ${term}`;

    try {
        const token = localStorage.getItem("token");
        //const teacherID = localStorage.getItem("teacher_id");
        const [studentsResponse, subjectsResponse, gradesResponse] = await Promise.all([
            fetch(`/api/classrooms/${classroomID}/students`, {
                method: 'GET', // Add the GET method
                headers: {
                    "Content-Type": "application/json",
                    'Authorization': `Bearer ${token}` // Add your authorization token here
                }
            }),
            fetch(`/api/classrooms/${classroomID}/subjects`, {
                method: 'GET', // Add the GET method
                headers: {
                    "Content-Type": "application/json",
                    'Authorization': `Bearer ${token}` // Add your authorization token here
                }
            }),
            fetch(`/api/classrooms/${classroomID}/terms/${termID}/grades`, {
                method: 'GET', // Add the GET method
                headers: {
                    "Content-Type": "application/json",
                    'Authorization': `Bearer ${token}` // Add your authorization token here
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
                const subjectContainer = document.createElement('div');
                subjectContainer.classList.add('subject-container');

                const subjectTitle = document.createElement('h2');
                subjectTitle.textContent = subject.name;
                subjectContainer.appendChild(subjectTitle);
                
                if (filteredGradeLabels.length > 0) {
                    const table = document.createElement('table');
                    table.id = `grades-grid-${subject.id}`;
                    table.classList.add('grades-grid');
                    subjectContainer.appendChild(table);
    
                    generateGradesGrid(table, students, filteredGradeLabels);
                } else {
                    const noGradeLabelsMessage = document.createElement('p');
                    noGradeLabelsMessage.textContent = "No grade labels added yet";
                    subjectContainer.appendChild(noGradeLabelsMessage);
                }

                gradesGridContainer.appendChild(subjectContainer);
            }
        );
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
        const gradesGridID = `grades-grid-${subject.id}`;
        const gradesGrid = document.getElementById(gradesGridID);
        if (gradesGrid) {
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
        } else {
            console.log(`No grades grid found for subject with ID ${subject.id}`);
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
    })
    .catch(error => {
        console.error('Error uploading grades:', error);
        alert('Failed to upload grades');
    });
}
*/

