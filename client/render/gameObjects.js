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
        this.sprite = sprite
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
            image(this.spriteIMG(), this.xblock * block_size, this.yblock * block_size, block_size, block_size);

            //fill('#f0d465');
            //rect(this.xblock * block_size, this.yblock * block_size, block_size, block_size);
            //rect(this.xpx, this.ypx, block_size, block_size);
        };
        //return this;
    }

    setPosition(x, y) {
        this.xblock = x
        this.yblock = y
    }

    spriteIMG() {
        return this.sprite
    }
}

class Pacman extends Entity {
    constructor(sprite) {
        super(sprite);
    }
}

class Ghost extends Entity {
    constructor(sprite) {
        super(sprite);
    }

    spriteIMG() {
        if (isScared) {
            return scared_image
        } else {
            return this.sprite
        }
    }
}