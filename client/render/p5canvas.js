//import {background, createCanvas, loadImage, windowHeight, windowWidth} from "p5/global";

let levelImg;
let pellets = [];
let entities = [];
let level;
loadJSON('./levels/level1.json', x => level = x)

let levelHeight, levelWidth;
let block_size;

function loadJSON(filePath, success, error) {
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            if (xhr.status === 200) {
                if (success)
                    success(JSON.parse(xhr.responseText));
            } else {
                if (error)
                    error(xhr);
            }
        }
    };
    xhr.open("GET", filePath, true);
    xhr.send();
}

function preload() {
    levelImg = loadImage('images/level.png')
    // level = levelStr.split('\n').map(str => str.split(' '))
    // print(data)
    print(level)
}

function calc_block_size() {
    levelHeight = windowHeight / 2
    levelWidth = windowHeight / 2 / levelImg.height * levelImg.width
    block_size = levelHeight / 25
}

function setup() {
    calc_block_size();
    createCanvas(levelWidth, levelHeight);
}

function windowResized() {
    calc_block_size();
    resizeCanvas(levelWidth, levelHeight);
}


function draw() {
    background(255);
    fill(255, 204, 0)
    noStroke()
    for (let i = 0; i < 25; i++) {
        for (let j = 0; j < 21; j++) {
            if (level[i][j] === '100') {
                rect(j * block_size, i * block_size, block_size, block_size)
            }
        }
    }

    // background(levelImg);
    // x++
}
