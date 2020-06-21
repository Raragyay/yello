//import {background, createCanvas, loadImage, windowHeight, windowWidth} from "p5/global";

let levelImg;
function preload(){
    levelImg=loadImage('images/level.png')
}

function setup() {
    createCanvas(windowHeight/2/levelImg.height*levelImg.width, windowHeight/2);
}

function windowResized(){
    createCanvas(windowHeight/2/levelImg.height*levelImg.width, windowHeight/2);
}


function draw() {
    background(levelImg);
    // x++
}
