// TODO REFACTOR INTO MOVEABLE OBJECT CLASS GENERIFY

import React from "react";
import {connect} from 'react-redux';
import {ReactComponent as PacmanSVG} from "../tileset/Pacman.svg";
import {loadPage} from "../actions";


class Pacman extends React.Component {

    componentDidMount() {
        this.props.dispatch(loadPage())
    }

    componentDidUpdate(prevProps, prevState, snapshot) {
        console.log(snapshot)
    }

    render() {
        const
            styles = {
                position: 'absolute',
                transform: `translate(${this.props.x}px,${this.props.y}px)`
            };
        return (
            <div style={styles}><PacmanSVG style={{width: 100, height: 100}}/></div>
        );
    }
}

const mapStateToProps = (state) => {
    return {
        x: state.pacmanState.x,
        y: state.pacmanState.y
    }
}

export default connect(mapStateToProps)(Pacman);
