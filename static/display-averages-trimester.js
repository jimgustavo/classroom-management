//  static/display-averages-trimester.js
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
            fetch(`/classroom/${classroomID}/averageswithfactors_trimesters?trimestre_1=0.7&sumativa_t1=0.3&trimestre_2=0.7&sumativa_t2=0.3&trimestre_3=0.7&sumativa_t3=0.3`, {
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

    const partialAverage3HeaderCell = headerRow.insertCell();
    partialAverage3HeaderCell.textContent = `Partial Average 3`;

    const finalAverageHeaderCell = headerRow.insertCell();
    finalAverageHeaderCell.textContent = `Final Average`;

    students.forEach((student, index) => {
        const row = gridElement.insertRow();

        const numberCell = row.insertCell();
        numberCell.textContent = index + 1;

        const nameCell = row.insertCell();
        nameCell.textContent = student.name;

        const studentAverage = averagesData.averages.find(a => a.student_id === student.id && a.subject_id === subjectID);
        if (studentAverage) {
            studentAverage.averages.forEach(termAverage => {
                const termAverageCell = row.insertCell();
                termAverageCell.contentEditable = false;

                const factorAverageCell = row.insertCell();
                factorAverageCell.contentEditable = false;

                // Update: Check if reinforcement is included and append (R.A.)
                if (termAverage.label === 'includes_reinforcement') {
                    termAverageCell.textContent = `${termAverage.average.toFixed(2)} (R.A.)`;
                } else {
                    termAverageCell.textContent = termAverage.average.toFixed(2);
                }
                
                factorAverageCell.textContent = termAverage.ave_factor.toFixed(2);
            });
            
            // Display partial averages
            const partialAverage1Cell = row.insertCell();
            partialAverage1Cell.textContent = studentAverage.partial_ave_1.toFixed(2);

            const partialAverage2Cell = row.insertCell();
            partialAverage2Cell.textContent = studentAverage.partial_ave_2.toFixed(2);

            const partialAverage3Cell = row.insertCell();
            partialAverage3Cell.textContent = studentAverage.partial_ave_3.toFixed(2);

            // Display final average
            const finalAverageCell = row.insertCell();
            finalAverageCell.textContent = studentAverage.term_ave.toFixed(2);
        }
    });
}
