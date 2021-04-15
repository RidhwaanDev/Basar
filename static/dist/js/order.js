const handleImageUpload = event => {
    alert("hello")
    const files = event.target.files
    const formData = new FormData()
    const name = document.getElementById("firstName").value
    const emailOrNumber = document.getElementById("emailOrNumber").value
    const madrasah = document.getElementById("madrasah").value
    
    formData.append('myFile', files[0])

    if (!name || !emailOrNumber || !madrasah){
        return
    }

    const orderData = {
        "name" : name,
        "contact" : emailOrNumber
        "madrasah" : madrasah
    }

    fetch('localhost:8000/upload', {
        method: 'POST',
        body: {
            formData,
            orderData
        }
    })
        .then(response => response.json())
        .then(data => {
            console.log(data.path)
        })
        .catch(error => {
            console.error(error)
        })
}
document.getElementsByClassName('form-group mt-4 mb-0')[0]
        .addEventListener('click', function (event) {
            alert("hello world")
        });


document.querySelector('#fileUpload').addEventListener('change', event => {
    alert("teAWD")
    handleImageUpload(event)
})

