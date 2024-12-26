// static/classroom-grades-display.js
document.addEventListener("DOMContentLoaded", async () => {
    const urlParams = new URLSearchParams(window.location.search);
    const classroomID = urlParams.get('classroomID');
    const termID = urlParams.get('termID');
    const term = urlParams.get('term');
    const classroomName = urlParams.get('classroomName');
    
    document.getElementById("classroom-title").textContent = `Grades Grid of ${classroomName} - ${term}`;

    try {
        const [studentsResponse, subjectsResponse, gradesResponse, reinforcementGradesResponse] = await Promise.all([
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
            }),
            fetch(`/grade-labels/reinforcement/classroom/${classroomID}/term/${termID}`, {
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
        const reinforcementGrades = await reinforcementGradesResponse.json() || []; // Ensure it's an array

        console.log("Fetched Grades Data:", gradesData);
        console.log("Fetched Reinforcement Grades Data:", reinforcementGrades);
       
        const gradesGridContainer = document.getElementById("grades-grid-container");

        subjects.forEach(subject => {
            const filteredGradeLabels = subject.grade_labels.filter(label => label.term_id == termID);
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

                generateGradesGrid(table, students, filteredGradeLabels, gradesData, reinforcementGrades, subject.id, term);

                gradesGridContainer.appendChild(subjectContainer);
            }
        });
    } catch (error) {
        console.error('Error fetching students, subjects, or grades:', error);
    }
});

function generateGradesGrid(gridElement, students, gradeLabels, gradesData, reinforcementGrades, subjectID, term) {
    const headerRow = gridElement.insertRow();

    const numberHeaderCell = headerRow.insertCell();
    numberHeaderCell.textContent = 'Number';

    const nameHeaderCell = headerRow.insertCell();
    nameHeaderCell.textContent = 'Student Name';

    gradeLabels.forEach(label => {
        const headerCell = headerRow.insertCell();
        headerCell.textContent = label.label; // Display label name
    });

    // Create a set of unique reinforcement labels
    const reinforcementLabels = [...new Set(reinforcementGrades.map(rg => rg.label))];
    const reinforcementColumns = reinforcementLabels.map(label => {
        const headerCell = headerRow.insertCell();
        headerCell.textContent = label; // Display reinforcement label name
        return { label, cellIndex: headerCell.cellIndex };
    });

    const averageHeaderCell = headerRow.insertCell();
    averageHeaderCell.textContent = `${term}-average`;

    const includesReinforcementHeaderCell = headerRow.insertCell();
    includesReinforcementHeaderCell.textContent = 'Includes Reinforcement';

    students.forEach((student, index) => {
        const row = gridElement.insertRow();

        const numberCell = row.insertCell();
        numberCell.textContent = index + 1;

        const nameCell = row.insertCell();
        nameCell.textContent = student.name;

        let totalGrades = 0;
        let gradeCount = 0;
        let includesReinforcement = false;

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

        // Populate reinforcement grades into their respective columns
        reinforcementColumns.forEach(column => {
            const reinforcementCell = row.insertCell(column.cellIndex);
            reinforcementCell.contentEditable = false;

            const studentReinforcementGrades = reinforcementGrades.filter(rg => rg.student_id === student.id && rg.subject_id === subjectID && rg.label === column.label);
            if (studentReinforcementGrades.length > 0) {
                studentReinforcementGrades.forEach(rg => {
                    reinforcementCell.textContent = rg.grade;
                    totalGrades += parseFloat(rg.grade);
                    gradeCount++;
                    includesReinforcement = true;
                });
            }
        });

        // Calculate and add term average
        const averageCell = row.insertCell();
        const average = gradeCount > 0 ? (totalGrades / gradeCount).toFixed(2) : '0.00';
        averageCell.textContent = average;

        // Add includes reinforcement cell
        const includesReinforcementCell = row.insertCell();
        includesReinforcementCell.textContent = includesReinforcement ? 'Yes' : 'No';
    });
}

/*
document.addEventListener("DOMContentLoaded", async () => {
    const urlParams = new URLSearchParams(window.location.search);
    const classroomID = urlParams.get('classroomID');
    const termID = urlParams.get('termID');
    const term = urlParams.get('term');
    const classroomName = urlParams.get('classroomName');
    
    document.getElementById("classroom-title").textContent = `Grades Grid of ${classroomName} - ${term}`;

    try {
        const [studentsResponse, subjectsResponse, gradesResponse, reinforcementGradesResponse] = await Promise.all([
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
            }),
            fetch(`/grade-labels/reinforcement/classroom/${classroomID}/term/${termID}`, {
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
        const reinforcementGrades = await reinforcementGradesResponse.json() || []; // Ensure it's an array

        console.log("Fetched Grades Data:", gradesData);
        console.log("Fetched Reinforcement Grades Data:", reinforcementGrades);
       
        const gradesGridContainer = document.getElementById("grades-grid-container");

        subjects.forEach(subject => {
            const filteredGradeLabels = subject.grade_labels.filter(label => label.term_id == termID);
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

                generateGradesGrid(table, students, filteredGradeLabels, gradesData, reinforcementGrades, subject.id, term);

                gradesGridContainer.appendChild(subjectContainer);
            }
        });
    } catch (error) {
        console.error('Error fetching students, subjects, or grades:', error);
    }
});

function generateGradesGrid(gridElement, students, gradeLabels, gradesData, reinforcementGrades, subjectID, term) {
    const headerRow = gridElement.insertRow();

    const numberHeaderCell = headerRow.insertCell();
    numberHeaderCell.textContent = 'Number';

    const nameHeaderCell = headerRow.insertCell();
    nameHeaderCell.textContent = 'Student Name';

    gradeLabels.forEach(label => {
        const headerCell = headerRow.insertCell();
        headerCell.textContent = label.label; // Display label name
    });

    const averageHeaderCell = headerRow.insertCell();
    averageHeaderCell.textContent = `${term}-average`;

    const includesReinforcementHeaderCell = headerRow.insertCell();
    includesReinforcementHeaderCell.textContent = 'Includes Reinforcement';

    students.forEach((student, index) => {
        const row = gridElement.insertRow();

        const numberCell = row.insertCell();
        numberCell.textContent = index + 1;

        const nameCell = row.insertCell();
        nameCell.textContent = student.name;

        let totalGrades = 0;
        let gradeCount = 0;
        let includesReinforcement = false;

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

        // Aggregate reinforcement grades if available
        const studentReinforcementGrades = reinforcementGrades.filter(rg => rg.student_id === student.id && rg.subject_id === subjectID);
        if (studentReinforcementGrades.length > 0) {
            studentReinforcementGrades.forEach(rg => {
                totalGrades += parseFloat(rg.grade);
                gradeCount++;
                includesReinforcement = true;
            });
        }

        // Calculate and add term average
        const averageCell = row.insertCell();
        const average = gradeCount > 0 ? (totalGrades / gradeCount).toFixed(2) : '0.00';
        averageCell.textContent = average;

        // Add includes reinforcement cell
        const includesReinforcementCell = row.insertCell();
        includesReinforcementCell.textContent = includesReinforcement ? 'Yes' : 'No';
    });
}
*/
/*
document.addEventListener("DOMContentLoaded", async () => {
    const urlParams = new URLSearchParams(window.location.search);
    const classroomID = urlParams.get('classroomID');
    const termID = urlParams.get('termID');
    const term = urlParams.get('term');
    const classroomName = urlParams.get('classroomName');
    
    document.getElementById("classroom-title").textContent = `Grades Grid of ${classroomName} - ${term}`;

    try {
        const [studentsResponse, subjectsResponse, gradesResponse, reinforcementGradesResponse] = await Promise.all([
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
            }),
            fetch(`/grade-labels/reinforcement/classroom/${classroomID}/term/${termID}`, {
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
        const reinforcementGrades = await reinforcementGradesResponse.json() || []; // Ensure it's an array

        console.log("Fetched Grades Data:", gradesData);
        console.log("Fetched Reinforcement Grades Data:", reinforcementGrades);
       
        const gradesGridContainer = document.getElementById("grades-grid-container");

        subjects.forEach(subject => {
            const filteredGradeLabels = subject.grade_labels.filter(label => label.term_id == termID);
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

                generateGradesGrid(table, students, filteredGradeLabels, gradesData, reinforcementGrades, subject.id, term);

                gradesGridContainer.appendChild(subjectContainer);
            }
        });
    } catch (error) {
        console.error('Error fetching students, subjects, or grades:', error);
    }
});

function generateGradesGrid(gridElement, students, gradeLabels, gradesData, reinforcementGrades, subjectID, term) {
    const headerRow = gridElement.insertRow();

    const numberHeaderCell = headerRow.insertCell();
    numberHeaderCell.textContent = 'Number';

    const nameHeaderCell = headerRow.insertCell();
    nameHeaderCell.textContent = 'Student Name';

    gradeLabels.forEach(label => {
        const headerCell = headerRow.insertCell();
        headerCell.textContent = label.label; // Display label name
    });

    // Check if there are any reinforcement grades
    const hasReinforcementGrades = reinforcementGrades.length > 0;
    let reinforcementHeaderCell;

    if (hasReinforcementGrades) {
        reinforcementHeaderCell = headerRow.insertCell();
        reinforcementHeaderCell.textContent = 'Reinforcement Grade';
    }

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

        // Add reinforcement grade cell if applicable
        if (hasReinforcementGrades) {
            const reinforcementCell = row.insertCell();
            reinforcementCell.contentEditable = false;

            const studentReinforcementGrades = reinforcementGrades.find(rg => rg.student_id === student.id && rg.subject_id === subjectID);
            if (studentReinforcementGrades) {
                reinforcementCell.textContent = studentReinforcementGrades.grade;
                totalGrades += parseFloat(studentReinforcementGrades.grade);
                gradeCount++;
            }
        }

        // Calculate and add term average
        const averageCell = row.insertCell();
        const average = gradeCount > 0 ? (totalGrades / gradeCount).toFixed(2) : '0.00';
        averageCell.textContent = average;
    });
}
*/