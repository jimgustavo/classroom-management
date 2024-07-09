//  static/admin.js

document.addEventListener("DOMContentLoaded", () => {
    fetchClassrooms();
    fetchTeachers()
    fetchAcademicPeriods();
    fetchTerms(); 
    fetchTeacherData();
    setTeacherNameInNavBar();    // Set the teacher's name in the navigation bar

    const logoutBtn = document.getElementById("logout-btn");
    logoutBtn.addEventListener("click", logout);

    document.getElementById("academicPeriod-form").addEventListener("submit", createAcademicPeriod);

    document.getElementById("assignTerm-form").addEventListener("submit", assignTermToAcademicPeriod);

    document.getElementById("close-terms-modal").addEventListener("click", closeTermsModal);

    const termForm = document.getElementById("term-form");
    termForm.addEventListener("submit", createTerm);

    document.getElementById("teacher-search-bar").addEventListener("input", filterTeachers);
    document.getElementById("assign-role-btn").addEventListener("click", assignRoleToTeacher);
});

///////////////**********TEACHERS SECTION**************////////////////////////////

let allTeachers = [];

async function fetchTeachers() {
    const token = localStorage.getItem("token");
    try {
        const response = await fetch(`/admin/teachers`, {
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });
        const teachers = await response.json();
        allTeachers = teachers; // Store all teachers for filtering
        displayTeachers(teachers);
    } catch (error) {
        console.error("Error fetching terms:", error);
    }
}

function displayTeachers(teachers) {
    const teachersList = document.getElementById("teachers-list");
    teachersList.innerHTML = "";
    teachers.forEach(teacher => {
        const li = document.createElement("li");
        li.textContent = `ID: ${teacher.id} Email: ${teacher.email} Name: ${teacher.name} Role: ${teacher.role}`;

        // Create delete button
        const deleteBtn = document.createElement("button");
        deleteBtn.textContent = "Delete";
        deleteBtn.classList.add("delete-btn");

        // Attach click event listener to delete button
        deleteBtn.addEventListener("click", () => {
            deleteTeacher(teacher.id); // Assuming each term has an 'id' property
        });

        li.appendChild(deleteBtn);
        teachersList.appendChild(li);
    });
}

function filterTeachers() {
    const searchTerm = document.getElementById("teacher-search-bar").value.toLowerCase();
    const filteredTeachers = allTeachers.filter(teacher => 
        teacher.name.toLowerCase().includes(searchTerm) || 
        teacher.id.toString().includes(searchTerm) ||
        teacher.email.toLowerCase().includes(searchTerm)
    );
    displayTeachers(filteredTeachers);
}

async function assignRoleToTeacher() {
    const searchBar = document.getElementById("teacher-search-bar");
    const roleDropdown = document.getElementById("role-dropdown");

    const teacherID = searchBar.value.trim();
    const role = roleDropdown.value;

    if (!teacherID) {
        alert("Please enter a valid teacher ID.");
        return;
    }

    const token = localStorage.getItem("token");
    try {
        const response = await fetch(`/admin/teacher/${teacherID}/role/${role}`, {
            method: "PUT",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json"
            }
        });

        if (response.ok) {
            alert("Role assigned successfully");
            fetchTeachers(); // Refresh the teacher list
        } else {
            const error = await response.json();
            alert(`Error: ${error.message}`);
        }
    } catch (error) {
        console.error("Error assigning role:", error);
    }
}

async function deleteTeacher(teacherID) {
    const token = localStorage.getItem("token");
    try {
        const response = await fetch(`/admin/teacher/${teacherID}`, {
            method: "DELETE",
            headers: {
                "Authorization": `Bearer ${token}`
            },
        });
        if (response.ok) {
            // Remove term from UI
            const teacherItem = document.querySelector(`[data-teacher-id="${teacherID}"]`);
            if (teacherItem) {
                teacherItem.remove();
            }
            fetchTeachers(); // Refresh the teacher list
        } else {
            console.error("Failed to delete term:", response.statusText);
        }
    } catch (error) {
        console.error("Error deleting term:", error);
    }
}



// Function to set the teacher's name in the navigation bar
function setTeacherNameInNavBar() {
    const teacherName = localStorage.getItem("teacher_name");
    if (teacherName) {
        const teacherNavItem = document.getElementById("teacher-nav-item");
        const teacherDataH2 = document.getElementById("teacher-data-h2");
        teacherNavItem.querySelector("a").textContent = teacherName;
        teacherDataH2.textContent = `${teacherName} Data`;
    }
}

async function logout() {
    try {
        // Clear token from local storage
        localStorage.removeItem("token");
        localStorage.removeItem("teacher_id");
        alert("Logout successful!");
        // Redirect to login page or perform any other action as needed
        window.location.href = "index.html";
    } catch (error) {
        console.error("Error logging out:", error);
        alert("An error occurred during logout. Please try again later.");
    }
}

async function fetchTeacherData() {
    const token = localStorage.getItem("token");
    try {
        const response = await fetch(`/admin/teacherdata`, {
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });
        const teacherData = await response.json();
        displayTeacherData(teacherData);
    } catch (error) {
        console.error("Error fetching teacher data:", error);
    }
}

function displayTeacherData(teacherData) {
    const teacherDataContainer = document.getElementById("teacherdata-container");
    teacherDataContainer.innerHTML = "";

    const teacherDiv = document.createElement("div");
    teacherDiv.classList.add("teacher-data-entry");
    teacherDiv.dataset.teacherId = teacherData.id;

    const fields = [
        `School: ${teacherData.school}`,
        `School Year: ${teacherData.school_year}`,
        `School Hours: ${teacherData.school_hours}`,
        `Country: ${teacherData.country}`,
        `City: ${teacherData.city}`,
        `Full Name: ${teacherData.teacher_full_name}`,
        `Birthday: ${teacherData.teacher_birthday}`,
        `ID Number: ${teacherData.id_number}`,
        `Labor Dependency Relationship: ${teacherData.labor_dependency_relationship}`,
        `Institutional Email: ${teacherData.institutional_email}`,
        `Phone: ${teacherData.phone}`,
        `Principal: ${teacherData.principal}`,
        `Vice Principal: ${teacherData.vice_principal}`,
        `DECE: ${teacherData.dece}`,
        `Inspector: ${teacherData.inspector}`
    ];

    fields.forEach(field => {
        const p = document.createElement("p");
        p.textContent = field;
        teacherDiv.appendChild(p);
    });
    teacherDataContainer.appendChild(teacherDiv);
}

///////////////**********CLASSROOM SECTION**************////////////////////////////

// Update fetchClassrooms and createClassroom to include token in headers
async function fetchClassrooms() {
    const token = localStorage.getItem("token");
    try {
        const response = await fetch(`/admin/classrooms`, {
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });
        const classrooms = await response.json();
        displayClassrooms(classrooms);
    } catch (error) {
        console.error("Error fetching classrooms:", error);
    }
}

// Function to display classrooms as cards
function displayClassrooms(classrooms) {
    const classroomList = document.getElementById("classroom-container");
    classroomList.innerHTML = "";
    
    classrooms.forEach(classroom => {
        const card = document.createElement("div");
        card.classList.add("classroom-cards");
        card.dataset.classroomId = classroom.id;

        const idPara = document.createElement("p");
        idPara.textContent = `ID: ${classroom.id}`;
        
        const heading = document.createElement("h3");
        heading.textContent = `${classroom.name}`;
        
        // Create delete button for classroom
        const deleteBtn = document.createElement("button");
        deleteBtn.textContent = "Delete Classroom";
        deleteBtn.classList.add("delete-btn");
        deleteBtn.addEventListener("click", () => {
            deleteClassroom(classroom.id);
        });
           
        card.appendChild(idPara);
        card.appendChild(heading);
        card.appendChild(deleteBtn);
    
        classroomList.appendChild(card);
    });
}

///////////////*****ACADEMIC PERIOD SECTION********/////////////////////////

// Fetch and display academic periods
async function fetchAcademicPeriods() {
    const token = localStorage.getItem("token");
    try {
        const response = await fetch(`/admin/academic_periods`, {
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });
        const academicPeriods = await response.json();
        displayAcademicPeriods(academicPeriods);
    } catch (error) {
        console.error("Error fetching academic periods:", error);
    }
}

function displayAcademicPeriods(academicPeriods) {
    const academicPeriodsContainer = document.getElementById("academicPeriods-container");
    const academicPeriodSelect = document.getElementById("academicPeriod-select");
    academicPeriodsContainer.innerHTML = "";
    academicPeriodSelect.innerHTML = "";

    academicPeriods.forEach(period => {
        const card = document.createElement("div");
        card.classList.add("classroom-cards");
        card.dataset.academicPeriodId = period.id;

        const idPara = document.createElement("p");
        idPara.textContent = `ID: ${period.id}`;

        const name = document.createElement("h3");
        name.textContent = `${period.name}`;

         // Create terms button
         const termsBtn = document.createElement("button");
         termsBtn.textContent = "Terms";
         termsBtn.classList.add("terms-btn");
         termsBtn.addEventListener("click", () => {
             fetchTermsInAcademicPeriod(period.id);
         });

        // Create delete button
        const deleteBtn = document.createElement("button");
        deleteBtn.textContent = "Delete";
        deleteBtn.classList.add("delete-btn");
        deleteBtn.addEventListener("click", () => {
            deleteAcademicPeriod(period.id);
        });

        card.appendChild(idPara);
        card.appendChild(name);
        card.appendChild(termsBtn);
        card.appendChild(deleteBtn);

        academicPeriodsContainer.appendChild(card);

        // Add option to select dropdown
        const option = document.createElement("option");
        option.value = period.id;
        option.textContent = period.name;
        academicPeriodSelect.appendChild(option);
    });
}

// Fetch terms in an academic period
async function fetchTermsInAcademicPeriod(academicPeriodId) {
    const token = localStorage.getItem("token");
    try {
        const response = await fetch(`/admin/academic_periods/${academicPeriodId}/terms`, {
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });
        const terms = await response.json();
        displayTermsInModal(terms);
    } catch (error) {
        console.error("Error fetching terms in academic period:", error);
    }
}

// Display terms in modal
function displayTermsInModal(terms) {
    const modalBody = document.getElementById("terms-modal-body");
    modalBody.innerHTML = "";
    
    terms.forEach(term => {
        const termItem = document.createElement("p");
        termItem.textContent = `Term: ${term.name}`;
        modalBody.appendChild(termItem);
    });

    const modal = document.getElementById("termsModal");
    modal.style.display = "block";
}

// Close terms modal
function closeTermsModal() {
    const modal = document.getElementById("termsModal");
    modal.style.display = "none";
}

// Create a new academic period
async function createAcademicPeriod(event) {
    event.preventDefault();

    const academicPeriodName = document.getElementById("academicPeriod-name").value;
    const token = localStorage.getItem("token");

    try {
        const response = await fetch("/admin/academic_periods", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            },
            body: JSON.stringify({ name: academicPeriodName })
        });

        if (response.ok) {
            alert("Academic Period created successfully!");
            fetchAcademicPeriods(); // Refresh the list of academic periods
        } else {
            const errorData = await response.json();
            alert(`Failed to create academic period: ${errorData.message}`);
        }
    } catch (error) {
        console.error("Error creating academic period:", error);
        alert("An error occurred while creating the academic period. Please try again later.");
    }
}

// Delete an academic period
async function deleteAcademicPeriod(academicPeriodId) {
    const token = localStorage.getItem("token");

    try {
        const response = await fetch(`/admin/academic_periods/${academicPeriodId}`, {
            method: "DELETE",
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });

        if (response.ok) {
            alert("Academic Period deleted successfully!");
            fetchAcademicPeriods(); // Refresh the list of academic periods
        } else {
            console.error("Failed to delete academic period:", response.statusText);
        }
    } catch (error) {
        console.error("Error deleting academic period:", error);
    }
}

// Assign a term to an academic period
async function assignTermToAcademicPeriod(event) {
    event.preventDefault();

    const academicPeriodId = document.getElementById("academicPeriod-select").value;
    const termId = document.getElementById("term-select").value;
    const token = localStorage.getItem("token");

    try {
        const response = await fetch(`/admin/academic_periods/${academicPeriodId}/terms/${termId}`, {
            method: "POST",
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });

        if (response.ok) {
            alert("Term assigned to academic period successfully!");
        } else {
            console.error("Failed to assign term to academic period:", response.statusText);
        }
    } catch (error) {
        console.error("Error assigning term to academic period:", error);
    }
}

function displayTermsDropdown(terms) {
    const termSelect = document.getElementById("term-select");
    termSelect.innerHTML = "";

    terms.forEach(term => {
        const option = document.createElement("option");
        option.value = term.id;
        option.textContent = term.name;
        termSelect.appendChild(option);
    });
}

// Close modal when clicking outside of it
window.onclick = function(event) {
    const modal = document.getElementById("termsModal");
    if (event.target == modal) {
        modal.style.display = "none";
    }
}

///////////////**********TERM SECTION**************////////////////////////////

async function fetchTerms() {
    try {
        const token = localStorage.getItem("token");
        const response = await fetch(`/admin/terms`, {
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });
        const terms = await response.json();
        displayTerms(terms);
        displayTermsDropdown(terms);
    } catch (error) {
        console.error("Error fetching terms:", error);
    }
}

function displayTerms(terms) {
    const termList = document.getElementById("term-list");
    termList.innerHTML = "";
    terms.forEach(term => {
        const li = document.createElement("li");
        li.textContent = `Term: ${term.name}`;

        // Create delete button
        const deleteBtn = document.createElement("button");
        deleteBtn.textContent = "Delete";
        deleteBtn.classList.add("delete-btn");

        // Attach click event listener to delete button
        deleteBtn.addEventListener("click", () => {
            deleteTerm(term.id); // Assuming each term has an 'id' property
        });

        li.appendChild(deleteBtn);
        termList.appendChild(li);
    });
}

async function createTerm(event) {
    event.preventDefault();

    const termName = document.getElementById("term-name").value;
    const academicPeriodId = document.getElementById("academicPeriod-select").value;
    const token = localStorage.getItem("token");

    // Replace spaces with underscores in the term name
    const sanitizedTermName = termName.replace(/\s+/g, "_");

    try {
        const response = await fetch("/admin/terms", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            },
            body: JSON.stringify({ name: sanitizedTermName, academic_period_id: parseInt(academicPeriodId) })
        });
        if (response.ok) {
            alert("Term created successfully!");
            fetchTerms(); // Refresh the term list
        } else {
            const errorData = await response.json();
            alert(`Failed to create term: ${errorData.message}`);
        }
    } catch (error) {
        console.error("Error creating term:", error);
        alert("An error occurred while creating the term. Please try again later.");
    }
}

async function deleteTerm(termId) {
    try {
        const response = await fetch(`/admin/terms/${termId}`, {
            method: "DELETE",
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            },
        });
        if (response.ok) {
            // Remove term from UI
            const termItem = document.querySelector(`[data-term-id="${termId}"]`);
            if (termItem) {
                termItem.remove();
            }
            fetchTerms(); // Refresh the term list
        } else {
            console.error("Failed to delete term:", response.statusText);
        }
    } catch (error) {
        console.error("Error deleting term:", error);
    }
}

