function submitData() {
    let name = document.getElementById("Name-input")
        .value;
    let email = document
        .getElementById("Email-input")
        .value;
    let phone = document
        .getElementById("Phone-input")
        .value;
    let subject = document
        .getElementById("subject-input")
        .value;
    let message = document
        .getElementById("input-message")
        .value;

    if (name == "") {
        return alert("Nama harus diisi dong");
    } else if (email == "") {
        return alert("Email harus diisi dong");
    } else if (phone == "") {
        return alert("Phone harus diisi dong");
    } else if (subject == "") {
        return alert("Subject harus dipilih dong");
    } else if (message == "") {
        return alert("Message harus diisi dong");
    }

    let emailReceiver = "isnan344@gmail.com";

    let a = document.createElement("a");
    a.href = `mailto:${emailReceiver}?subject=${subject}&body=hello, nama saya ${name}, ${message}. kamu bisa menghubungi saya di nomer ${phone}, terima kasih.`;
    a.click();
}