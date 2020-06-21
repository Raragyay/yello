
// Initialize a sound classifier method with SpeechCommands18w model.
let classifier;
const options = { probabilityThreshold: 0.8 };
// Two variables to hold the label and confidence of the result
let label;
let confidence;

async function setup() {
  classifier = await ml5.soundClassifier("SpeechCommands18w", options);
  // Create 'label' and 'confidence' div to hold results

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
  // Display error in the console
  if (error) {
    console.error(error);
  }
  // The results are in an array ordered by confidence.
  console.log(results);
  // Show the first label and confidence
  label.textContent = "Label: " + results[0].label;
  confidence.textContent = "Confidence: " + results[0].confidence.toFixed(4);
}
