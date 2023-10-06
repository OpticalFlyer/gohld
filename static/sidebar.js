// Toggle sidebar
document.getElementById("toggle-sidebar").addEventListener("click", function() {
    var sidebar = document.getElementById("sidebar");
    
    if (sidebar.style.width === "0px" || sidebar.style.width === "") {
        sidebar.style.width = "250px"; // Or your desired width
    } else {
        sidebar.style.width = "0px";
    }
});
