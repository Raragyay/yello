//import {background, createCanvas, loadImage, windowHeight, windowWidth} from "p5/global";

let pellets = [];
let entities = [];
let level, isWall;
loadJSON('./levels/level1.json', x => {
    level = x.levelData
    isWall = level.map(x => x.map(tile => tile === '100'))
})
let levelHeight, levelWidth;
let block_size;
let player1 = new Pacman;

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
    // levelImg = loadImage('images/level.png')
    // level = levelStr.split('\n').map(str => str.split(' '))
    // print(data)
    // print(isWall)
}

function calc_block_size() {
    levelHeight = windowHeight / 2
    levelWidth = windowHeight / 2 / level.length * level[0].length
    block_size = levelHeight / level.length
}

function setup() {
    calc_block_size();
    var canvas = createCanvas(levelWidth, levelHeight);
    // canvas.parent('sketch-div')
    player1 = Pacman();
}

function windowResized() {
    calc_block_size();
    resizeCanvas(levelWidth, levelHeight);
}


function drawLevel() {
    fill(255, 204, 0)
    noStroke()
    for (let i = 0; i < isWall.length; i++) {
        for (let j = 0; j < isWall[0].length; j++) {
            if (isWall[i][j]) {
                let border_checks = []
                border_checks[0] = (j > 0 && isWall[i][j - 1])
                border_checks[1] = (i > 0 && isWall[i - 1][j])
                border_checks[2] = (j < isWall[0].length - 1 && isWall[i][j + 1])
                border_checks[3] = (i < isWall.length - 1 && isWall[i + 1][j])
                let corner_checks = []
                for (let k = 0; k < 4; k++) {
                    corner_checks[k] = !((border_checks[k] || border_checks[(k + 1) % 4])) * 10
                }
                rect(j * block_size, i * block_size, block_size, block_size, ...corner_checks)
            }
        }
    }
}

function draw() {
    background(255);
    drawLevel();
    player1.update();
    player1.show();
}


function Pacman() {
    this.x = 0;
    this.y = 0;
    this.xspeed = 1;
    this.yspeed = 0;

    this.update = function () {

        if (isWall[this.x + this.xspeed][this.y + this.yspeed]) {
            this.xspeed = 0;
            this.yspeed = 0;
        } else {
            this.x = this.x + this.xspeed;
            this.y = this.y + this.yspeed;
        }

        this.show = function () {
            fill(0);
            rect(this.x, this.y, 100, 100);
        }
    }
    return this;
}

function updateCommand(newCmd) {
    command = newCmd;
    switch (newCmd) {
        case 'up':
            player1.xspeed = 0;
            player1.yspeed = -1;
            break;
        case 'down':
            player1.xspeed = 0;
            player1.yspeed = 1;
            break;
        case 'left':
            player1.xspeed = -1;
            player1.yspeed = 0;
        case 'right':
            player1.xspeed = 1;
            player1.yspeed = 0;
    }
}

// Initialize a sound classifier method with SpeechCommands18w model.
let classifier;
const options = {probabilityThreshold: 0.8};
// Two variables to hold the label and confidence of the result
let label;
let confidence;

let command;

let wordToCmd = {};
/*{
  red: 'left',
  yellow: 'up',
  green: 'right',
  blue: 'down'
};*/

let cmdToWord = {
    left: 'red',
    up: 'yellow',
    right: 'green',
    down: 'blue'
}


async function setupstt() {
    classifier = await ml5.soundClassifier(
        "https://storage.googleapis.com/tm-model/RoRt49x-Z/model.json",
        options
    );

    updateDicts('red', 'yellow', 'green', 'blue');

    // Create 'label' and 'confidence' div to hold results,, delete eventually

    label = document.createElement("DIV");
    label.textContent = "label ...";
    confidence = document.createElement("DIV");
    confidence.textContent = "Confidence ...";

    document.body.appendChild(label);
    document.body.appendChild(confidence);
    // Classify the sound from microphone in real time
    classifier.classify(gotResult);

}

setupstt();
console.log("ml5 version:", ml5.version);

// A function to run when we get any errors and the results
function gotResult(error, results) {
    // for debug
    if (error) {
        console.error(error);
    }

    let wordIn = results[0].label;
    label.textContent = "Label: " + wordIn;
    confidence.textContent = "Confidence: " + results[0].confidence.toFixed(4);

    updateCommand(wordToCmd[wordIn]);
}


function updateDicts(newLeft, newUp, newDown, newRight) {
    cmdToWord.left = newLeft;
    cmdToWord.up = newUp;
    cmdToWord.right = newRight;
    cmdToWord.down = newDown;

    //update reverse dict
    for (var key in cmdToWord) {
        if (cmdToWord.hasOwnProperty(key)) {
            wordToCmd[cmdToWord[key]] = key;
        }
    }
}
