//  static/reports.js

// Fetch and populate classrooms and terms on page load
document.addEventListener("DOMContentLoaded", () => {
    fetchClassrooms();
    fetchTerms();
    fetchSubjects();
});

async function fetchClassrooms() {
    try {
        const token = localStorage.getItem("token");
        const teacherID = localStorage.getItem("teacher_id");
       
        const response = await fetch(`/api/classrooms/teacher/${teacherID}`, {
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const classrooms = await response.json();
        populateClassroomDropdown(classrooms); // Populate dropdown with classrooms
        populateReportClassroomDropdown(classrooms)
        populateClassroomDropdownForAverages(classrooms)
    } catch (error) {
        console.error("Error fetching classrooms:", error);
    }
}

let subjects = []; // Declare a global variable to store subjects

async function fetchSubjects() {
    try {
        const token = localStorage.getItem("token");
        const teacherID = localStorage.getItem("teacher_id");

        const response = await fetch(`/api/subjects/teacher/${teacherID}`, {
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        subjects = await response.json();
        console.log("Subjects fetched:", subjects); // Debug log
    } catch (error) {
        console.error("Error fetching subjects:", error);
    }
}

async function fetchTerms() {
    try {
        const token = localStorage.getItem("token");
        const academicPeriodID = localStorage.getItem("academic_period");
        console.log("academic period id:", academicPeriodID);
        const response = await fetch(`/api/academic_periods/${academicPeriodID}/terms`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const terms = await response.json();
        populateTermDropdown(terms); // Populate dropdown with terms
        populateReportTermDropdown(terms)
    } catch (error) {
        console.error("Error fetching terms:", error);
    }
}

function populateClassroomDropdown(classrooms) {
    const classroomDropdown = document.getElementById("classroom-dropdown");

    classrooms.forEach(classroom => {
        const option = document.createElement("option");
        option.value = classroom.id;
        option.textContent = classroom.name;
        classroomDropdown.appendChild(option);
    });

    classroomDropdown.addEventListener("change", handleClassroomOrTermChange);
}

// Populate classroom dropdown  for report generation in final average
function populateClassroomDropdownForAverages(classrooms) {
    const finalAverageDropdown = document.getElementById("final-average-report-dropdown");

    classrooms.forEach(classroom => {
        const option = document.createElement("option");
        option.value = classroom.id;
        option.textContent = classroom.name;
        finalAverageDropdown.appendChild(option);
    });
}

function populateTermDropdown(terms) {
    const termDropdown = document.getElementById("term-dropdown");

    terms.forEach(term => {
        const option = document.createElement("option");
        option.value = term.id;
        option.textContent = term.name;
        termDropdown.appendChild(option);
    });

    termDropdown.addEventListener("change", handleClassroomOrTermChange);
}

async function handleClassroomOrTermChange() {
    const classroomDropdown = document.getElementById("classroom-dropdown");
    const termDropdown = document.getElementById("term-dropdown");

    const classroomID = classroomDropdown.value;
    const termID = termDropdown.value;

    if (classroomID && termID) {
        try {
            const [studentsResponse, gradesResponse] = await Promise.all([
                fetch(`/api/classrooms/${classroomID}/students`, {
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

            if (!studentsResponse.ok || !gradesResponse.ok) {
                throw new Error(`HTTP error! status: ${studentsResponse.status}, ${gradesResponse.status}`);
            }

            const students = await studentsResponse.json();
            const gradesData = await gradesResponse.json();

            displayStudentsWithLowGrades(students, gradesData, classroomID, termID);

        } catch (error) {
            console.error("Error fetching students or grades:", error);
        }
    }
}

function displayStudentsWithLowGrades(students, gradesData, classroomID, termID) {
    const token = localStorage.getItem("token");
    const teacherID = localStorage.getItem("teacher_id");
    const role = localStorage.getItem("role");
    
    const lowGradesContainer = document.getElementById("low-grades-container");
    lowGradesContainer.innerHTML = "";

    // Retrieve the selected term's name
    const termDropdown = document.getElementById("term-dropdown");
    const selectedTermName = termDropdown.options[termDropdown.selectedIndex].text;

    students.forEach(student => {
        const studentGrades = gradesData.grades.find(g => g.student_id === student.id);

        if (studentGrades) {
            // Match the term by its name instead of ID
            const termGrades = studentGrades.terms.find(t => t.term === selectedTermName);
            const lowGrades = termGrades ? termGrades.grades.filter(g => g.grade < 7) : [];

            if (lowGrades.length > 0) {
                const studentDiv = document.createElement("div");
                studentDiv.classList.add("student-entry");

                const studentName = document.createElement("h3");
                studentName.textContent = student.name;
                studentDiv.appendChild(studentName);

                lowGrades.forEach(grade => {
                    const gradeDiv = document.createElement("div");
                    gradeDiv.classList.add("grade-entry");

                    // Find the subject name using the subject_id
                    const subject = subjects.find(subject => subject.id === studentGrades.subject_id);
                    const subjectName = subject ? subject.name : 'Unknown Subject';
 
                    const gradeLabel = document.createElement("p");
                    gradeLabel.textContent = `Subject: ${subjectName}, Grade: ${grade.grade}`;
                    gradeDiv.appendChild(gradeLabel);

                    const addReinforcementButton = document.createElement("button");
                    addReinforcementButton.textContent = "Añadir Refuerzo";
                    addReinforcementButton.addEventListener("click", () => {
                        openReinforcementModal(student.id, classroomID, subject.id, termID);
                    });
                    gradeDiv.appendChild(addReinforcementButton);

                    const generateMinuteButton = document.createElement("button");
                    generateMinuteButton.textContent = "Acta de Compromiso";
                    generateMinuteButton.disabled = role !== "proteacher"; // Disable by default, enable for proteacher
                    generateMinuteButton.addEventListener("click", async () => {
                        if (role === "proteacher") {
                            try {
                                const response = await fetch(`/proteacher/pdfminute/teacher/${teacherID}/classroom/${classroomID}/student/${student.id}`, {
                                    method: 'GET',
                                    headers: {
                                        "Authorization": `Bearer ${token}`
                                    }
                                });

                                if (!response.ok) {
                                    throw new Error(`HTTP error! status: ${response.status}`);
                                }

                                const blob = await response.blob();
                                const url = window.URL.createObjectURL(blob);
                                const a = document.createElement('a');
                                a.style.display = 'none';
                                a.href = url;
                                a.download = `acta-de-compromiso.pdf`;
                                document.body.appendChild(a);
                                a.click();
                                window.URL.revokeObjectURL(url);
                            } catch (error) {
                                console.error("Error generating report:", error);
                            }
                        } else {
                            alert("Only pro teachers can generate reports.");
                        }
                    });
                    gradeDiv.appendChild(generateMinuteButton);

                    studentDiv.appendChild(gradeDiv);
                });

                lowGradesContainer.appendChild(studentDiv);
            }
        }
    });
}

function openReinforcementModal(studentID, classroomID, subjectID, termID) {
    const modal = document.getElementById("reinforcementModal");
    const student = document.getElementById("reinforcement-student");
    const classroom = document.getElementById("reinforcement-classroom");
    const subject = document.getElementById("reinforcement-subject");
    const term = document.getElementById("reinforcement-term");
    
    student.value = studentID;
    classroom.value = classroomID;
    subject.value = subjectID;
    term.value = termID;

    //console.log(`studentID: ${studentID}, subjectID: ${subjectID}, classroomID: ${classroomID} and termID: ${termID}`);

    modal.style.display = "block";
}

function closeModal() {
    const modal = document.getElementById("reinforcementModal");
    modal.style.display = "none";
}

document.getElementById("reinforcementForm").addEventListener("submit", async (event) => {
    event.preventDefault();

    const studentID = parseInt(document.getElementById("reinforcement-student").value);
    const classroomID = parseInt(document.getElementById("reinforcement-classroom").value);
    const subjectID = parseInt(document.getElementById("reinforcement-subject").value);
    const termID = parseInt(document.getElementById("reinforcement-term").value);
    const gradeLabel = document.getElementById("reinforcement-grade-label").value;
    const date = document.getElementById("reinforcement-date").value;
    const skill = document.getElementById("reinforcement-skill").value;
    const teacherID = parseInt(localStorage.getItem("teacher_id"));
    const grade = parseFloat(document.getElementById("reinforcement-grade").value);

    const bodyTest = JSON.stringify({ 
                student_id: studentID, 
                classroom_id: classroomID, 
                subject_id: subjectID, 
                term_id: termID, 
                label: gradeLabel, 
                date: date, 
                skill: skill, 
                teacher_id: teacherID, 
                grade: grade 
            })
    console.log(`reinforcement bodyTest: ${bodyTest}`);

    try {
        const response = await fetch("/api/grade-labels/reinforcement", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            },
            body: JSON.stringify({ 
                student_id: studentID, 
                classroom_id: classroomID, 
                subject_id: subjectID, 
                term_id: termID, 
                label: gradeLabel, 
                date: date, 
                skill: skill, 
                teacher_id: teacherID, 
                grade: grade 
            })
        });

        if (response.ok) {
            alert("Reinforcement grade label added successfully!");
            closeModal();
        } else {
            const errorData = await response.json();
            alert(`Error: ${errorData.error}`);
        }
    } catch (error) {
        console.error("Error adding reinforcement grade label:", error);
    }
});

// Populate classroom dropdown for report generation
function populateReportClassroomDropdown(classrooms) {
    const classroomReportDropdown = document.getElementById("classroom-report-dropdown");

    classrooms.forEach(classroom => {
        const option = document.createElement("option");
        option.value = classroom.id;
        option.textContent = classroom.name;
        classroomReportDropdown.appendChild(option);
    });
}

// Populate term dropdown for report generation
function populateReportTermDropdown(terms) {
    const termReportDropdown = document.getElementById("term-report-dropdown");

    terms.forEach(term => {
        const option = document.createElement("option");
        option.value = term.id;
        option.textContent = term.name;
        termReportDropdown.appendChild(option);
    });
}

// Generate XLSX report
document.getElementById("generate-xlsx-report").addEventListener("click", async () => {
    const classroomReportDropdown = document.getElementById("classroom-report-dropdown");
    const termReportDropdown = document.getElementById("term-report-dropdown");

    const classroomID = classroomReportDropdown.value;
    const termID = termReportDropdown.value;

    const teacherID = localStorage.getItem("teacher_id");
    const academicPeriodID = localStorage.getItem("academic_period");

    if (classroomID && termID) {
        try {
            const response = await fetch(`/xlsx-report/teachers/${teacherID}/classrooms/${classroomID}/academicPeriod/${academicPeriodID}/terms/${termID}`, {
                method: 'GET',
                headers: {
                    "Authorization": `Bearer ${localStorage.getItem("token")}`
                }
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.style.display = 'none';
            a.href = url;
            a.download = `reporte-de-calificaciones.xlsx`;
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
        } catch (error) {
            console.error("Error generating XLSX report:", error);
        }
    } else {
        alert("Please select both a classroom and a term.");
    }
});

// Generate Average XLSX report
document.getElementById("generate-xlsx-average-report").addEventListener("click", async () => {
    const averageReportDropdown = document.getElementById("final-average-report-dropdown");

    const classroomID = averageReportDropdown.value;

    const teacherID = localStorage.getItem("teacher_id");
    const academicPeriodID = localStorage.getItem("academic_period");

    if (classroomID) {
        try {
            const response = await fetch(`/xlsx-average/teachers/${teacherID}/classrooms/${classroomID}/academicPeriod/${academicPeriodID}`, {
                method: 'GET',
                headers: {
                    "Authorization": `Bearer ${localStorage.getItem("token")}`
                }
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.style.display = 'none';
            a.href = url;
            a.download = `promedios-finales.xlsx`;
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
        } catch (error) {
            console.error("Error generating XLSX report:", error);
        }
    } else {
        alert("Please select a classroom");
    }
});



