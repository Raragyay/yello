export const RIGHT = 'RIGHT';
export const LEFT = 'LEFT';
export const UP = 'UP';
export const DOWN = 'DOWN';
export const PACMAN_MOVE = 'PACMAN_MOVE';
// Call these actions with the
export const moveRight = () => (
    {
        type: PACMAN_MOVE,
        direction: RIGHT
    })
export const moveLeft = () => (
    {
        type: PACMAN_MOVE,
        direction: LEFT
    })
export const moveUp = () => (
    {
        type: PACMAN_MOVE,
        direction: UP
    })
export const moveDown = () => (
    {
        type: PACMAN_MOVE,
        direction: DOWN
    })
export const loadPage = () => {
    return function (dispatch, getState) {
        window.addEventListener('keydown', (e) => {
            switch (e.code) {
                case "ArrowRight":
                    e.preventDefault();
                    dispatch(moveRight());
                    break;

                case "ArrowLeft":
                    e.preventDefault();
                    dispatch(moveLeft());
                    break;

                case "ArrowDown":
                    e.preventDefault();
                    dispatch(moveDown());
                    break;

                case "ArrowUp":
                    e.preventDefault();
                    dispatch(moveUp());
                    break;

                default:
            }
        })
    }
};
