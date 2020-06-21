// Initialize a sound classifier method with SpeechCommands18w model.
let classifier;
const options = { probabilityThreshold: 0.8 };
// Two variables to hold the label and confidence of the result
let label;
let confidence;

import {updateCommand} from 'p5canvas.js';

export let command;

export let wordToCmd = {};
/*{
  red: 'left',
  yellow: 'up',
  green: 'right',
  blue: 'down'
};*/

export let cmdToWord = {
  left: 'red',
  up: 'yellow',
  right: 'green',
  down: 'blue'
}


async function setup() {
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

setup();
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


function updateDicts(newLeft, newUp, newDown, newRight){
  cmdToWord.left = newLeft;
  cmdToWord.up = newUp;
  cmdToWord.right = newRight;
  cmdToWord.down = newDown;

  //update reverse dict
  for(var key in cmdToWord) {
    if(cmdToWord.hasOwnProperty(key)){
      wordToCmd[cmdToWord[key]] = key;
    }
  }
}