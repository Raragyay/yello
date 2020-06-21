//import {background, createCanvas, loadImage, windowHeight, windowWidth} from "p5/global";
// Initialize a sound classifier method with SpeechCommands18w model.
let classifier;
const options = {probabilityThreshold: 0.7};
// Two variables to hold the label and confidence of the result
let label;
let confidence;
let command;
let player1;
let canvas;
let gameActive = false;

var canvasDiv = document.getElementById('canvas-div')

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


let pellets = [];
let entities = [];
let level, isWall;
let levelHeight, levelWidth;
let block_size;

entities.push(new Entity())


function loadTXT(filePath, success, error) {
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            if (xhr.status === 200) {
                if (success)
                    success(xhr.responseText);
                console.log('successfully loaded level')
            } else {
                if (error)
                    error(xhr);
            }
        }
    };
    xhr.open("GET", filePath, true);
    xhr.send();
}

let player_image, blinky_image, inky_image, pinky_image, clyde_image

function preload() {
    loadTXT('./levels/level2.txt', load_level)
    player_image = loadImage('images/player_image.png')
    blinky_image = loadImage('images/blinky_image.png')
    inky_image = loadImage('images/inky_image.png')
    pinky_image = loadImage('images/pinky_image.png')
    clyde_image = loadImage('images/clyde_image.png')
}

const waitForLevelToBeDefined = () => {
    return new Promise((resolve) => {
        if (level === undefined) {
            return setTimeout(() => resolve(waitForLevelToBeDefined), 1000);
        } else {
            return resolve()
        }
    })
}

function calc_block_size(cb = function () {
}) {
    waitForLevelToBeDefined().then(() => {
        levelHeight = windowHeight / 2
        levelWidth = windowHeight / 2 / level.length * level[0].length
        block_size = levelHeight / level.length
        //console.log('bscalced')
    }).then(cb);
}


async function setup() {
    calc_block_size(() => {
        canvas = createCanvas(levelWidth, levelHeight);
        console.log("createdCanvas");
        canvas.parent('canvas-div')
        canvas.center();
        canvas.style.position = "relative";
        canvas.style('z-index', "-3");
    });
}

async function windowResized() {
    calc_block_size(() => {
        resizeCanvas(levelWidth, levelHeight);
        canvas.center();
    });
    //console.log("windowresized");
}


function drawLevel() {
    fill('#57d4ef')
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
    pellets.forEach(pellet => {
            pellet.draw()
        }
    )
}

function draw() {
    clear();
    // if (!gameActive) {
    //     return
    // }
    if (level === undefined || player1 === undefined) {
        return
    }
    drawLevel();
    player1.update();
    player1.show();
}

class Pellet {
    constructor(x_, y_) {
        this.x = x_
        this.y = y_
    }

    draw() {
        fill("yellow")
        circle(this.x * block_size + block_size / 2, this.y * block_size + block_size / 2, this.pellet_size())
    }

    pellet_size() {
        return block_size / 4
    }
}

class BigPellet extends Pellet {
    pellet_size() {
        return block_size / 2
    }
}

class Entity {
    constructor(sprite) {
        this.sprite=sprite
        this.xblock = 10;
        this.yblock = 2;
        //this.xpx = this.xblock * block_size;
        //this.ypx = this.yblock * block_size;
        this.xspeed = 0;
        this.yspeed = 0;

        this.update = function () {
            if (isWall[this.yblock + this.yspeed][this.xblock + this.xspeed]) {
                this.xspeed = 0
                this.yspeed = 0
            } else {
                this.xblock += this.xspeed;
                this.yblock += this.yspeed;
            }
        };
        this.show = function () {
            fill('#f0d465');
            rect(this.xblock * block_size, this.yblock * block_size, block_size, block_size);
            //rect(this.xpx, this.ypx, block_size, block_size);
        };
        //return this;
    }
}

class Player extends Entity {
    constructor(sprite) {
        super(sprite);
    }
}


function updateCommand(newCmd) {
    console.log("newcmd");
    if (!gameActive) {
        return
    }
    command = newCmd;
    switch (newCmd) {
        case 'up':
            player1.xspeed = 0;
            player1.yspeed = -1;
            sendSocketMessage("PONG UPDATE-DIR U")
            break;
        case 'down':
            player1.xspeed = 0;
            player1.yspeed = 1;
            sendSocketMessage("PONG UPDATE-DIR D")
            break;
        case 'left':
            player1.xspeed = -1;
            player1.yspeed = 0;
            sendSocketMessage("PONG UPDATE-DIR L")
            break;
        case 'right':
            player1.xspeed = 1;
            player1.yspeed = 0;
            sendSocketMessage("PONG UPDATE-DIR R")
            break;
        default:
            break;
    }
}


async function setupstt() {
    classifier = await ml5.soundClassifier(
        "https://storage.googleapis.com/tm-model/RoRt49x-Z/model.json",
        options
    );

    updateDicts('red', 'yellow', 'green', 'blue');

    // Create 'label' and 'confidence' div to hold results,, delete eventually

    label = document.createElement("DIV");
    label.textContent = "Command:";
    confidence = document.createElement("DIV");
    confidence.textContent = "Confidence:";
    confidence.classList.add("labels");
    label.classList.add("labels");

    document.body.appendChild(label);
    //document.body.appendChild(confidence);
    // Classify the sound from microphone in real time
    classifier.classify(gotResult);

}


console.log("ml5 versijlaflaon:", ml5.version);

// A function to run when we get any errors and the results
function gotResult(error, results) {
    // for debug
    if (error) {
        console.error(error);
    }

    let wordIn = results[0].label;
    label.textContent = "Command: " + wordIn;
    confidence.textContent = "Confidence: " + results[0].confidence.toFixed(4);

    updateCommand(wordToCmd[wordIn]);
}


function updateDicts(newLeft, newUp, newRight, newDown) {
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

setupstt();

//arrow key navigation
document.onkeydown = checkKey;

function checkKey(e) {

    e = e || window.event;
    if (e.keyCode == '38') {
        // up arrow
        updateCommand('up');
    } else if (e.keyCode == '40') {
        // down arrow
        updateCommand('down');
    } else if (e.keyCode == '37') {
        // left arrow
        updateCommand('left');
    } else if (e.keyCode == '39') {
        // right arrow
        updateCommand('right');
    }

}
