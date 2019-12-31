async function newBoatTextUpdate() {

    const boatName = await newBoatValidation(boatInput.value);
    console.log(boatName);

};




async function getBoatList() {
    try {
        const url = `http://localhost:8080/api?boat`;
        const result = await fetch(url);
        const data = await result.json();
        return data;

    } catch(error) {
        console.log(error);

    }
}

// on load - get the data from the api
document.onload = function() {

    const boats = await getBoatList();

    const element = document.getElementById('newBoatInput');
    element.addEventListener('input', update);
    


};

// add a checker to validate the string 


// when valid - allow user to submit

