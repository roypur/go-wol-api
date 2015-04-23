$(document).ready(
function()
{
    $("#btn").click(function()
    {
        wake($("#host").val(), $("#user").val(),$("#pass").val(),function(response){console.log(response)});
    })
});
