(function(){
				let loader = document.getElementById("loader");
				loader.style.display = 'none'
												document.getElementById('uploadButton').addEventListener('click', uploadFile);
// change button to loading bar
// document.getElementById('uploadButton').addEventListener('click', testFileDownload);

function uploadFile(){
				let headers = new Headers();
				headers.append('Origin','https://scanpdf.app:8000');

				let file = document.getElementById("fileUpload").files[0];
				if(!file){
								alert("no file")
								return
				}
				let formData = new FormData();

				formData.append("myFile", file);
				//upload file
				loader.style.display = 'block';
				fetch('https://scanpdf.app:8000/upload', {method: "POST", body: formData, headers:headers})
								.then(response =>  {
												let obj = response.json().then(data => {
																let id = data['Id'];
																console.log(id)
																periodCheck(id)
												})
								})
				.catch(err => {
								loader.style.display ='none'
								alert(err);
				});
}
function periodCheck(id){
				let headers = new Headers();
				headers.append('Origin','https://scanpdf');
				let url = 'https://scanpdf.app:8000/checkTicket?' + new URLSearchParams({ id : id})

				let timer = setInterval(function(){
								console.log("checking if job is done")
								fetch(url, {method: "GET", headers:headers})
												.then(response => {
																let obj = response.json().then(data => {
																				console.log("response 2" + data.toString());
																				let status = data['Status'];
																				if(status == 2){
																								alert("Done")
																								clearInterval(timer)
																								testFileDownload(data['FileName'])
																				}
																})
												})
				}, 3000)
}

function testFileDownload(fileName){
				console.log(`getting ${fileName}`)
				axios({
								url: `https://scanpdf.app:8000/${fileName}`,
								method: 'GET',
								responseType: 'blob', // important
				}).then((response) => {
								const url = window.URL.createObjectURL(new Blob([response.data]));
								const link = document.createElement('a');
								link.href = url;
								link.setAttribute('download', 'result.txt');
								document.body.appendChild(link);
								link.click();
				});
}
}())
