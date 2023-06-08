const promise = new Promise ((resolve, reject) => {
    const xhr = new XMLHttpRequest();
    xhr.open('GET', 'https://api.npoint.io/cf15ab5d6f21dcab5069', true);
    xhr.onload = () => {
        if (xhr.status === 200) {
            resolve(JSON.parse(xhr.response));
        } else {
            reject("Error loading data");
        }
    };

    xhr.onerror = () => {
        reject("Network error");
    };

    xhr.send();
});

function starMaker(rating) {

    if (rating === 1) {
        return `<i class="fa-solid fa-star"><i class="fa-regular fa-star"></i><i class="fa-regular fa-star"></i><i class="fa-regular fa-star"></i><i class="fa-regular fa-star"></i></i>`;
    } else if (rating === 2) {
        return `<i class="fa-solid fa-star"></i><i class="fa-solid fa-star"></i><i class="fa-regular fa-star"></i><i class="fa-regular fa-star"></i><i class="fa-regular fa-star"></i>`;
    } else if (rating === 3) {
        return `<i class="fa-solid fa-star"></i><i class="fa-solid fa-star"></i><i class="fa-solid fa-star"></i><i class="fa-regular fa-star"></i><i class="fa-regular fa-star"></i>`;
    } else if (rating === 4) {
        return `<i class="fa-solid fa-star"></i><i class="fa-solid fa-star"></i><i class="fa-solid fa-star"></i><i class="fa-solid fa-star"></i><i class="fa-regular fa-star"></i>`;
    } else if ( rating == 5 ) {
        return `<i class="fa-solid fa-star"></i><i class="fa-solid fa-star"></i><i class="fa-solid fa-star"></i><i class="fa-solid fa-star"></i><i class="fa-solid fa-star"></i>`;
    }
}

async function allTesti() {
    const response = await promise;

    let testimonialHTML = "";
    response.forEach(item => {
        testimonialHTML += `<div class="col-6 col-md-3">
                                <div class="card p-3 border-0 shadow">
                                    <img src="${item.image}" alt="">
                                    <p class="mt-3">${item.quotes}</p>
                                    <p class="text-end">- ${item.author}</p>
                                    <p class="text-end">${starMaker(item.rating)}</p>
                                </div>
                            </div>`
    });

    document.getElementById('testi-container').innerHTML = testimonialHTML;
}

allTesti();

async function filterTesti(rating) {
    const response = await promise 

    const testiFiltered = response.filter(item => item.rating === rating);

    let testimonialHTML = "";

    if (testiFiltered.length === 0) {
        testimonialHTML = `<h1 class="text-center">404 Data Not Found<h1>`
    } else {
        testiFiltered.forEach(item => 
        testimonialHTML += `<div class="col-6 col-md-3">
                                <div class="card p-3 border-0 shadow">
                                    <img src="${item.image}" alt="">
                                    <p class="mt-3">${item.quotes}</p>
                                    <p class="text-end">- ${item.author}</p>
                                    <p class="text-end">${starMaker(item.rating)}</p>
                                </div>
                            </div>`
        )};
    
        document.getElementById('testi-container').innerHTML = testimonialHTML;
}