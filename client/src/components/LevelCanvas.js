import React from "react";
import level from "../tileset/level.png"

class LevelCanvas extends React.Component {
    constructor(props) {
        super(props);
        this.canvasRef = React.createRef();
        this.levelRef = React.createRef();
    }

    componentDidMount() {
        const ctx = this.canvasRef.current.getContext('2d');
        this.refs.image.onLoad = () => {
            ctx.drawImage(this.refs.image, 0, 0);
        }
    }

    render() {
        return (
            <div>
                <canvas ref={this.canvasRef}/>
                <img ref="image" src={level} style={{display: 'none'}}/>
            </div>
        )
    }
}

export default LevelCanvas
