//  static/display-averages.js
document.addEventListener("DOMContentLoaded", async () => {
    const urlParams = new URLSearchParams(window.location.search);
    const classroomID = urlParams.get('classroomID');
    const classroomName = urlParams.get('classroomName');
    
    document.getElementById("classroom-title").textContent = `Average Grid of ${classroomName}`;

    try {
        const [studentsResponse, subjectsResponse, averagesResponse] = await Promise.all([
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
            fetch(`/classroom/${classroomID}/averageswithreinforcement?bimestre_1=0.8&sumativa_1=0.2&bimestre_2=0.8&sumativa_2=0.2`, {
                method: 'GET', 
                headers: {
                    "Content-Type": "application/json",
                }
            }),
        ]);

        const students = await studentsResponse.json();
        const subjects = await subjectsResponse.json();
        const averagesData = await averagesResponse.json();

        console.log("Fetched Students:", students);
        console.log("Fetched Subjects:", subjects);
        console.log("Fetched Average Data:", averagesData);
       
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

            generateGradesGrid(table, students, averagesData, subject.id);

            gradesGridContainer.appendChild(subjectContainer);
        });
    } catch (error) {
        console.error('Error fetching students, subjects, or averages:', error);
    }
});

function generateGradesGrid(gridElement, students, averagesData, subjectID) {
    const headerRow = gridElement.insertRow();

    const numberHeaderCell = headerRow.insertCell();
    numberHeaderCell.textContent = 'Number';

    const nameHeaderCell = headerRow.insertCell();
    nameHeaderCell.textContent = 'Student Name';

    // Assuming terms are consistent across subjects
    const terms = averagesData.averages[0].averages.map(avg => avg.term);
    
    terms.forEach(term => {
        const termHeaderCell = headerRow.insertCell();
        termHeaderCell.textContent = term;
        
        const averageFactorHeaderCell = headerRow.insertCell();
        averageFactorHeaderCell.textContent = `%`;
    });

    const averageHeaderCell = headerRow.insertCell();
    averageHeaderCell.textContent = `Final Average`;

    const includesReinforcementHeaderCell = headerRow.insertCell();
    includesReinforcementHeaderCell.textContent = 'Includes Reinforcement';

    students.forEach((student, index) => {
        const row = gridElement.insertRow();

        const numberCell = row.insertCell();
        numberCell.textContent = index + 1;

        const nameCell = row.insertCell();
        nameCell.textContent = student.name;

        let totalAverage = 0;
        let termCount = 0;
        let includesReinforcement = false;

        terms.forEach(term => {
            const termAverageCell = row.insertCell();
            termAverageCell.contentEditable = false;

            const factorAverageCell = row.insertCell();
            factorAverageCell.contentEditable = false;

            const studentAverage = averagesData.averages.find(a => a.student_id === student.id && a.subject_id === subjectID);
            if (studentAverage) {
                const termAverage = studentAverage.averages.find(t => t.term === term);
                if (termAverage) {
                    if (termAverage.label === 'includes_reinforcement') {
                        termAverageCell.textContent = `${termAverage.average.toFixed(2)} (R.A.)`;
                        includesReinforcement = true;
                    } else {
                        termAverageCell.textContent = termAverage.average.toFixed(2);
                    }
                    factorAverageCell.textContent = termAverage.ave_factor.toFixed(2);
                    
                    totalAverage += parseFloat(termAverage.ave_factor);
                    termCount++;
                } else {
                    termAverageCell.textContent = 'N/A';
                    factorAverageCell.textContent = 'N/A';
                }
            } else {
                termAverageCell.textContent = 'N/A';
                factorAverageCell.textContent = 'N/A';
            }
        });

        const finalAverageCell = row.insertCell();
        const finalAverage = termCount > 0 ? (totalAverage / 2).toFixed(2) : '0.00';
        finalAverageCell.textContent = finalAverage;

        const includesReinforcementCell = row.insertCell();
        includesReinforcementCell.textContent = includesReinforcement ? 'Yes' : 'No';
    });
}

/*
document.addEventListener("DOMContentLoaded", async () => {
    const urlParams = new URLSearchParams(window.location.search);
    const classroomID = urlParams.get('classroomID');
    const classroomName = urlParams.get('classroomName');
    
    document.getElementById("classroom-title").textContent = `Average Grid of ${classroomName}`;

    try {
        const [studentsResponse, subjectsResponse, averagesResponse] = await Promise.all([
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
            fetch(`/classroom/${classroomID}/averageswithreinforcement?bimestre_1=0.8&sumativa_1=0.2&bimestre_2=0.8&sumativa_2=0.2`, {
                method: 'GET', 
                headers: {
                    "Content-Type": "application/json",
                }
            }),
        ]);

        const students = await studentsResponse.json();
        const subjects = await subjectsResponse.json();
        const averagesData = await averagesResponse.json();

        console.log("Fetched Students:", students);
        console.log("Fetched Subjects:", subjects);
        console.log("Fetched Average Data:", averagesData);
       
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

            generateGradesGrid(table, students, averagesData, subject.id);

            gradesGridContainer.appendChild(subjectContainer);
        });
    } catch (error) {
        console.error('Error fetching students, subjects, or averages:', error);
    }
});

function generateGradesGrid(gridElement, students, averagesData, subjectID) {
    const headerRow = gridElement.insertRow();

    const numberHeaderCell = headerRow.insertCell();
    numberHeaderCell.textContent = 'Number';

    const nameHeaderCell = headerRow.insertCell();
    nameHeaderCell.textContent = 'Student Name';

    // Assuming terms are consistent across subjects
    const terms = averagesData.averages[0].averages.map(avg => avg.term);
    
    terms.forEach(term => {
        const termHeaderCell = headerRow.insertCell();
        termHeaderCell.textContent = term;
        
        const averageFactorHeaderCell = headerRow.insertCell();
        averageFactorHeaderCell.textContent = `Average-%`;

        const reinforcementHeaderCell = headerRow.insertCell();
        reinforcementHeaderCell.textContent = `Reinforcement`;
    });

    const averageHeaderCell = headerRow.insertCell();
    averageHeaderCell.textContent = `Final Average`;

    const includesReinforcementHeaderCell = headerRow.insertCell();
    includesReinforcementHeaderCell.textContent = 'Includes Reinforcement';

    students.forEach((student, index) => {
        const row = gridElement.insertRow();

        const numberCell = row.insertCell();
        numberCell.textContent = index + 1;

        const nameCell = row.insertCell();
        nameCell.textContent = student.name;

        let totalAverage = 0;
        let termCount = 0;
        let includesReinforcement = false;

        terms.forEach(term => {
            const termAverageCell = row.insertCell();
            termAverageCell.contentEditable = false;

            const factorAverageCell = row.insertCell();
            factorAverageCell.contentEditable = false;

            const reinforcementCell = row.insertCell();
            reinforcementCell.contentEditable = false;

            const studentAverage = averagesData.averages.find(a => a.student_id === student.id && a.subject_id === subjectID);
            if (studentAverage) {
                const termAverage = studentAverage.averages.find(t => t.term === term);
                if (termAverage) {
                    termAverageCell.textContent = termAverage.average.toFixed(2);
                    factorAverageCell.textContent = termAverage.ave_factor.toFixed(2);

                    if (termAverage.label === 'includes_reinforcement') {
                        reinforcementCell.textContent = termAverage.average.toFixed(2);
                        includesReinforcement = true;
                    } else {
                        reinforcementCell.textContent = 'N/A';
                    }
                    
                    totalAverage += parseFloat(termAverage.ave_factor);
                    termCount++;
                } else {
                    termAverageCell.textContent = 'N/A';
                    factorAverageCell.textContent = 'N/A';
                    reinforcementCell.textContent = 'N/A';
                }
            } else {
                termAverageCell.textContent = 'N/A';
                factorAverageCell.textContent = 'N/A';
                reinforcementCell.textContent = 'N/A';
            }
        });

        const finalAverageCell = row.insertCell();
        const finalAverage = termCount > 0 ? (totalAverage / 2).toFixed(2) : '0.00';
        finalAverageCell.textContent = finalAverage;

        const includesReinforcementCell = row.insertCell();
        includesReinforcementCell.textContent = includesReinforcement ? 'Yes' : 'No';
    });
}
*/
/*
function generateGradesGrid(gridElement, students, averagesData, subjectID) {
    const headerRow = gridElement.insertRow();

    const numberHeaderCell = headerRow.insertCell();
    numberHeaderCell.textContent = 'Number';

    const nameHeaderCell = headerRow.insertCell();
    nameHeaderCell.textContent = 'Student Name';

    // Assuming terms are consistent across subjects
    const terms = averagesData.averages[0].averages.map(avg => avg.term);
    
    terms.forEach(term => {
        const termHeaderCell = headerRow.insertCell();
        termHeaderCell.textContent = term;
         // Add average-factor header
        const averageFactorHeaderCell = headerRow.insertCell();
        averageFactorHeaderCell.textContent = `Average-%`;
    });

    // Add final-average header
    const averageHeaderCell = headerRow.insertCell();
    averageHeaderCell.textContent = `Final Average`;

    students.forEach((student, index) => {
        const row = gridElement.insertRow();

        const numberCell = row.insertCell();
        numberCell.textContent = index + 1;

        const nameCell = row.insertCell();
        nameCell.textContent = student.name;

        let totalAverage = 0;
        let termCount = 0;

        terms.forEach(term => {
            const termAverageCell = row.insertCell();
            termAverageCell.contentEditable = false;

            const factorAverageCell = row.insertCell();
            factorAverageCell.contentEditable = false;

            // Find the corresponding average for this student, subject, and term
            const studentAverage = averagesData.averages.find(a => a.student_id === student.id && a.subject_id === subjectID);
            if (studentAverage) {
                const termAverage = studentAverage.averages.find(t => t.term === term);
                console.log("termAverage:", termAverage);
                if (termAverage) {
                    termAverageCell.textContent = termAverage.average.toFixed(2);
                    factorAverageCell.textContent = termAverage.ave_factor.toFixed(2);
                    totalAverage += parseFloat(termAverage.ave_factor);
                    termCount++;
                } else {
                    termAverageCell.textContent = 'N/A';
                }
            } else {
                termAverageCell.textContent = 'N/A';
            }
        });

        // Calculate and add final average
        const finalAverageCell = row.insertCell();
        const finalAverage = termCount > 0 ? (totalAverage/2).toFixed(2) : '0.00';
        finalAverageCell.textContent = finalAverage;
    });
}
*/
/*
document.addEventListener("DOMContentLoaded", async () => {
    const urlParams = new URLSearchParams(window.location.search);
    const classroomID = urlParams.get('classroomID');
    const classroomName = urlParams.get('classroomName');

    document.getElementById("classroom-title").textContent = `Average Grid of ${classroomName}`;

    try {
        const [studentsResponse, subjectsResponse, averagesResponse] = await Promise.all([
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
            fetch(`/classroom/${classroomID}/averageswithfactors?bimestre_1=0.8&sumativa_1=0.2&bimestre_2=0.8&sumativa_2=0.2`, {
                method: 'GET',
                headers: {
                    "Content-Type": "application/json",
                }
            }),
        ]);

        const students = await studentsResponse.json();
        const subjects = await subjectsResponse.json();
        const averagesData = await averagesResponse.json();

        console.log("Fetched Students:", students);
        console.log("Fetched Subjects:", subjects);
        console.log("Fetched Average Data:", averagesData);

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

            generateGradesGrid(table, students, averagesData, subject.id);

            gradesGridContainer.appendChild(subjectContainer);
        });
    } catch (error) {
        console.error('Error fetching students, subjects, or averages:', error);
    }
});

function generateGradesGrid(gridElement, students, averagesData, subjectID) {
    const headerRow = gridElement.insertRow();

    const numberHeaderCell = headerRow.insertCell();
    numberHeaderCell.textContent = 'Number';

    const nameHeaderCell = headerRow.insertCell();
    nameHeaderCell.textContent = 'Student Name';

    // Assuming terms are consistent across subjects
    const terms = averagesData.averages[0].averages.map(avg => avg.term);

    terms.forEach(term => {
        const termHeaderCell = headerRow.insertCell();
        termHeaderCell.textContent = term;
        // Add average-factor header
        const averageFactorHeaderCell = headerRow.insertCell();
        averageFactorHeaderCell.textContent = `%`;
    });

    // Add partial and final average headers
    const partialAverage1HeaderCell = headerRow.insertCell();
    partialAverage1HeaderCell.textContent = `Partial Average 1`;

    const partialAverage2HeaderCell = headerRow.insertCell();
    partialAverage2HeaderCell.textContent = `Partial Average 2`;

    const finalAverageHeaderCell = headerRow.insertCell();
    finalAverageHeaderCell.textContent = `Final Average`;

    students.forEach((student, index) => {
        const row = gridElement.insertRow();

        const numberCell = row.insertCell();
        numberCell.textContent = index + 1;

        const nameCell = row.insertCell();
        nameCell.textContent = student.name;

        let totalAverage = 0;
        let termCount = 0;
        let partialAverage1 = 0;
        let partialAverage2 = 0;

        const studentAverage = averagesData.averages.find(a => a.student_id === student.id && a.subject_id === subjectID);
        if (studentAverage) {
            studentAverage.averages.forEach((termAverage, termIndex) => {
                const termAverageCell = row.insertCell();
                termAverageCell.contentEditable = false;
                termAverageCell.textContent = termAverage.average.toFixed(2);

                const factorAverageCell = row.insertCell();
                factorAverageCell.contentEditable = false;
                factorAverageCell.textContent = termAverage.ave_factor.toFixed(2);

                if (termIndex < 2) {
                    partialAverage1 += termAverage.ave_factor;
                } else {
                    partialAverage2 += termAverage.ave_factor;
                }

                totalAverage += termAverage.ave_factor;
                termCount++;
            });
        }

        // Add partial averages
        const partialAverage1Cell = row.insertCell();
        partialAverage1Cell.textContent = partialAverage1.toFixed(2);

        const partialAverage2Cell = row.insertCell();
        partialAverage2Cell.textContent = partialAverage2.toFixed(2);

        // Calculate and add final average
        const finalAverageCell = row.insertCell();
        const finalAverage = termCount > 0 ? (totalAverage / 2).toFixed(2) : '0.00';
        finalAverageCell.textContent = finalAverage;
    });
}
*/
/*
document.addEventListener("DOMContentLoaded", async () => {
    const urlParams = new URLSearchParams(window.location.search);
    const classroomID = urlParams.get('classroomID');
    const classroomName = urlParams.get('classroomName');
    
    document.getElementById("classroom-title").textContent = `Average Grid of ${classroomName}`;

    try {
        const [studentsResponse, subjectsResponse, averagesResponse] = await Promise.all([
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
            fetch(`/classroom/${classroomID}/averageswithfactors?bimestre_1=0.8&sumativa_1=0.2&bimestre_2=0.8&sumativa_2=0.2`, {
                method: 'GET', 
                headers: {
                    "Content-Type": "application/json",
                }
            }),
        ]);

        const students = await studentsResponse.json();
        const subjects = await subjectsResponse.json();
        const averagesData = await averagesResponse.json();

        console.log("Fetched Students:", students);
        console.log("Fetched Subjects:", subjects);
        console.log("Fetched Average Data:", averagesData);
       
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

                generateGradesGrid(table, students, averagesData, subject.id);

                gradesGridContainer.appendChild(subjectContainer);
        });
    } catch (error) {
        console.error('Error fetching students, subjects, or averages:', error);
    }
});

function generateGradesGrid(gridElement, students, averagesData, subjectID) {
    const headerRow = gridElement.insertRow();

    const numberHeaderCell = headerRow.insertCell();
    numberHeaderCell.textContent = 'Number';

    const nameHeaderCell = headerRow.insertCell();
    nameHeaderCell.textContent = 'Student Name';

    // Assuming terms are consistent across subjects
    const terms = averagesData.averages[0].averages.map(avg => avg.term);
    
    terms.forEach(term => {
        const termHeaderCell = headerRow.insertCell();
        termHeaderCell.textContent = term;
         // Add average-factor header
        const averageFactorHeaderCell = headerRow.insertCell();
        averageFactorHeaderCell.textContent = `Average-%`;
    });

    // Add final-average header
    const averageHeaderCell = headerRow.insertCell();
    averageHeaderCell.textContent = `Final Average`;

    students.forEach((student, index) => {
        const row = gridElement.insertRow();

        const numberCell = row.insertCell();
        numberCell.textContent = index + 1;

        const nameCell = row.insertCell();
        nameCell.textContent = student.name;

        let totalAverage = 0;
        let termCount = 0;

        terms.forEach(term => {
            const termAverageCell = row.insertCell();
            termAverageCell.contentEditable = false;

            const factorAverageCell = row.insertCell();
            factorAverageCell.contentEditable = false;

            // Find the corresponding average for this student, subject, and term
            const studentAverage = averagesData.averages.find(a => a.student_id === student.id && a.subject_id === subjectID);
            if (studentAverage) {
                const termAverage = studentAverage.averages.find(t => t.term === term);
                console.log("termAverage:", termAverage);
                if (termAverage) {
                    termAverageCell.textContent = termAverage.average.toFixed(2);
                    factorAverageCell.textContent = termAverage.ave_factor.toFixed(2);
                    totalAverage += parseFloat(termAverage.ave_factor);
                    termCount++;
                } else {
                    termAverageCell.textContent = 'N/A';
                }
            } else {
                termAverageCell.textContent = 'N/A';
            }
        });

        // Calculate and add final average
        const finalAverageCell = row.insertCell();
        const finalAverage = termCount > 0 ? (totalAverage/2).toFixed(2) : '0.00';
        finalAverageCell.textContent = finalAverage;
    });
}
*/