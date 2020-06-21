import {combineReducers} from "redux";
import {DOWN, LEFT, PACMAN_MOVE, RIGHT, UP} from "../actions";

const blockUnit = 5;//TODO refactor to constants file

function pacmanState(state = {}, action) { //TODO REPLACE WITH SERVER GIVING INPUT
    switch (action.type) {
        case PACMAN_MOVE:
            switch (action.direction) {
                case RIGHT:
                    return Object.assign({}, state, {
                        ...state,
                        x: (state.x === undefined ? 0 : state.x) + blockUnit
                    });
                case LEFT:
                    return Object.assign({}, state, {
                        ...state,
                        x: (state.x === undefined ? 0 : state.x) - blockUnit
                    });
                case UP:
                    return Object.assign({}, state, {
                        ...state,
                        y: (state.y === undefined ? 0 : state.y) - blockUnit
                    });
                case DOWN:
                    return Object.assign({}, state, {
                        ...state,
                        y: (state.y === undefined ? 0 : state.y) + blockUnit
                    });
                default:
                    return state;
            }
        default:
            return state;
    }
}

const PacmanGame = combineReducers({
    pacmanState
});
export default PacmanGame;
