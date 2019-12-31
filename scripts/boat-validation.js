
const boatInput = document.getElementById('newBoatInput');

boatInput.addEventListener('input', update);

async function newBoatValidation() {

    try {
        const url = `http://localhost:8080/api/?boat`;
        const result = await fetch(url);
        const data = await result.json();
        return data;

    } catch(error) {
        console.log(error);

    }
};

async function newBoatTextUpdate() {

    const boatName = await newBoatValidation(boatInput.value);
    console.log(boatName);

};


function update() {
    newBoatTextUpdate();
}