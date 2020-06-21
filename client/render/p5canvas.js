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
/*loadJSON('./levels/level1.json', x => {
    level = x.levelData
    isWall = level.map(x => x.map(tile => tile === '100'))
})*/
let levelHeight, levelWidth;
let block_size;


function loadJSON(filePath, success, error) {
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            if (xhr.status === 200) {
                if (success)
                    success(JSON.parse(xhr.responseText));
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

function preload() {
    // levelImg = loadImage('images/level.png')
    // level = levelStr.split('\n').map(str => str.split(' '))
    loadJSON('./levels/level1.json', x => {
      level = x.levelData
      isWall = level.map(x => x.map(tile => tile === '100'))
      console.log(isWall);
  })
}

async function calc_block_size() {
    while (level === undefined) {
        await new Promise(resolve => {
            setTimeout(function () {
                resolve();
            }, 1000)
        })
    }
    levelHeight = windowHeight / 2
    levelWidth = windowHeight / 2 / level.length * level[0].length
    block_size = levelHeight / level.length
    //console.log('bscalced')
    return;
}

async function setup() {
    await calc_block_size();
    canvas = createCanvas(levelWidth, levelHeight);
    canvas.center();
    console.log("createdCanvas");
    // canvas.parent('sketch-div')
    player1 = new Player();
    console.log("setup");
    
    canvas.style.position = "relative";
    canvas.style('z-index', "-3");
}

async function windowResized() {
    await calc_block_size();
    resizeCanvas(levelWidth, levelHeight);
    canvas.center();
    //console.log("windowresized");
}


function drawLevel() {

  fill('#57D4EF')
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
   // console.log("drewLevel");
}

function draw() {
    clear();
    drawLevel();
    player1.update();
    player1.show();
}


class Player {
  constructor() {
    this.xblock = 10;
    this.yblock = 2;
    //this.xpx = this.xblock * block_size;
    //this.ypx = this.yblock * block_size;
    this.xspeed = 0;
    this.yspeed = 0;

    this.update = function () {
      if (isWall[this.yblock + this.yspeed][this.xblock + this.xspeed]) {
        this.xspeed = 0//0//this.x + this.xspeed*-0.1;
        this.yspeed = 0//0//this.y + this.yspeed*-0.1;
      }
      else {
        this.xblock += this.xspeed;
        this.yblock += this.yspeed;
        //this.xblock = Math.floor(this.xpx/block_size + block_size/2);
        //this.yblock = Math.floor(this.ypx/block_size + block_size/2);
        //console.log(this.xblock, this.yblock);
        //this.xpx = this.xpx + this.xspeed*0.1;
        //this.ypx = this.ypx + this.yspeed*0.1;
      }
    };
    this.show = function () {
      fill('#f0d465');
      rect(this.xblock*block_size, this.yblock*block_size, block_size, block_size);
      //rect(this.xpx, this.ypx, block_size, block_size);
    };
    //return this;
  }
}

function updateCommand(newCmd) {
    console.log("newcmd");
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
            break;
        case 'right':
            player1.xspeed = 1;
            player1.yspeed = 0;
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

    document.body.appendChild(label);
    document.body.appendChild(confidence);
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
    }
    else if (e.keyCode == '40') {
        // down arrow
        updateCommand('down');
    }
    else if (e.keyCode == '37') {
       // left arrow
       updateCommand('left');
    }
    else if (e.keyCode == '39') {
       // right arrow
       updateCommand('right');
    }

}
