// Toggle sidebar
document.getElementById("toggle-sidebar").addEventListener("click", function() {
    document.getElementById("users-sidebar").style.width = "0px";

    var sidebar = document.getElementById("sidebar");
    
    if (sidebar.style.width === "0px" || sidebar.style.width === "") {
        sidebar.style.width = "250px"; // Or your desired width
    } else {
        sidebar.style.width = "0px";
    }
});

// Toggle users sidebar
document.getElementById("users-button").addEventListener("click", function() {
    document.getElementById("sidebar").style.width = "0px";

    var sidebar = document.getElementById("users-sidebar");
    
    if (sidebar.style.width === "0px" || sidebar.style.width === "") {
        sidebar.style.width = "250px"; // Or your desired width
    } else {
        sidebar.style.width = "0px";
    }
});

document.getElementById("create-user-button").addEventListener("click", function() {
    document.getElementById("createUserModal").style.display = "block";
});