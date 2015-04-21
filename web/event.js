$(document).ready(
function()
{
    $("#btn").click(function()
    {
        console.log("click");
        
        console.log($("#user").val());
        
        sessionStorage.setItem('user', $("#user").val());
        
        window.sessionStorage.setItem('pass', $("#pass").val());
        
        
        wake("m1", $("#user").val(),$("#pass").val(),function(response){console.log(response)});
        
        
    })
});
