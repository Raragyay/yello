import isWall from './p5canvas.js';

function setup(){
  createCanvas(600, 600);
  player1 = new Pacman();
}

function draw(){
  background('#255');
  player1.update();
  player1.show();
}

function Pacman(){
  this.x = 0;
  this.y = 0;
  this.xspeed = 1;
  this.yspeed = 0;

  this.update = function() {
    /*if (isWall[this.x + this.xspeed][this.y + this.yspeed]){
      this.xspeed = 0;
      this.yspeed = 0;
    }
    else{*/
      this.x = this.x + this.xspeed;
      this.y = this.y + this.yspeed;
    //}

  this.show = function() {
    fill(0);
    rect(this.x, this.y, 100, 100);
  }
  }
}

