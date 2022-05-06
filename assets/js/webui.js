$("#formBootstrap").submit(function (event) {
    event.preventDefault();

    $("#buttonBootstrap")
        .prop("disabled", true)
        .attr("class", "btn btn-secondary btn-lg")
        .empty()
        .append('<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>\n')
        .append("Bootstrapping...");

    var postData = {
        inputHostname: $("#inputHostname").val(),
        inputHostgroup: $("#inputHostgroup").val(),
        inputGitRepo: $("#inputGitRepo").val(),
        inputVaultPassword: $("#inputVaultPassword").val()
    },
        postUrl = $(this).attr("action");

    $.post(postUrl, postData)
        .done(function (data) {
            console.log(data);
            $("#buttonBootstrap")
                .empty()
                .attr("class", "btn btn-success btn-lg")
                .text("Bootstrap done!");
        })
        .fail(function (data) {
            if (data.status === 0) {
                console.log("Error: Connection refused");
            } else {
                console.log("Error " + data.status + ': ' + data.statusText);
                console.log(data.responseJSON);
            }
            $("#buttonBootstrap")
                .empty()
                .attr("class", "btn btn-danger btn-lg")
                .text("Bootstrap failed!");
        });
});
