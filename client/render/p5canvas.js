import {command} from "./stt.js"
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
    var canvas=createCanvas(levelWidth, levelHeight);
    canvas.parent('sketch-div')
    player1 = new Pacman();
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

function isWall(x, y){
    return isWall[y][x];
}

function Pacman(){
    this.x = 0;
    this.y = 0;
    this.xspeed = 1;
    this.yspeed = 0;

    this.update = function() {
      
      if (isWall[this.x + this.xspeed][this.y + this.yspeed]){
        this.xspeed = 0;
        this.yspeed = 0;
      }
      else{
        this.x = this.x + this.xspeed;
        this.y = this.y + this.yspeed;
      }
  
    this.show = function() {
      fill(0);
      rect(this.x, this.y, 100, 100);
    }
    }
  }

function updateCommand(newCmd){
    command = newCmd;
    switch(newCmd){
      case up:
        player1.xspeed = 0;
        player1.yspeed = -1;
        break;
      case down:
        player1.xspeed = 0;
        player1.yspeed = 1;
        break;
      case left:
        player1.xspeed = -1;
        player1.yspeed = 0;
      case right:
        player1.xspeed = 1;
        player1.yspeed = 0;
    }
  }