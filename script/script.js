

const chartInstances = {};

const ignoredClasses = [
  'NPC',
  'Warmage',
  'Monk',
];

const dataUrl = 'data/allTheData';

const backgroundColors = [
  'rgba(255, 99, 132, 0.3)',  // Red
  'rgba(54, 162, 235, 0.3)',  // Blue
  'rgba(255, 235, 59, 0.3)',  // Bright yellow
  'rgba(0, 0, 0, 0.3)',       // Black
  'rgba(32, 201, 151, 0.3)',  // Greenish teal
  'rgba(153, 102, 255, 0.3)', // Purple
  'rgba(255, 159, 64, 0.3)',  // Orange
  'rgba(103, 58, 183, 0.3)',  // Deep purple
  'rgba(0, 150, 136, 0.3)',   // Dark teal
  'rgba(233, 30, 99, 0.3)',   // Pink
  'rgba(244, 67, 54, 0.3)',   // Deep orange
  'rgba(255, 206, 86, 0.3)',  // Yellow
  'rgba(75, 192, 192, 0.3)',  // Teal
];

const borderColors = [
  'rgba(255, 99, 132, 1)',
  'rgba(54, 162, 235, 1)',
  'rgba(255, 235, 59, 1)',
  'rgba(0, 0, 0, 1)',
  'rgba(32, 201, 151, 1)',
  'rgba(153, 102, 255, 1)',
  'rgba(255, 159, 64, 1)',
  'rgba(103, 58, 183, 1)',
  'rgba(0, 150, 136, 1)',
  'rgba(233, 30, 99, 1)',
  'rgba(244, 67, 54, 1)',
  'rgba(255, 206, 86, 1)',
  'rgba(75, 192, 192, 1)',
];

function getColor(path) {
  const pathToClassMap = {
    '1': 'color1',
    'albert': 'color1',
    '2': 'color2',
    'fähre': 'color2',
    '3': 'color3',
    'joel': 'color3',
    '4': 'color4',
    'kristoffer': 'color4',
    '5': 'color5',
    'markus': 'color5',
    '6': 'color6',
    'robban': 'color6',
    '7': 'color7',
    'vilijam': 'color7',
    '8': 'color8',
    '9': 'color9',
    '10': 'color10',
    '11': 'color11',
    '12': 'color12',
    '13': 'color13',
  };
  return pathToClassMap[path] || 'defaultClass'; // Use 'defaultClass' for any unspecified paths
}

function showImage(src) {
  // Set the source of the popup image
  document.getElementById('popupImage').src = src;
  // Display the popup
  document.getElementById('imagePopup').style.display = 'flex';
}

function closePopup() {
  // Hide the popup
  document.getElementById('imagePopup').style.display = 'none';
  // Remove the source of the popup image to not unnecessarily keep the image in memory
  document.getElementById('popupImage').src = '';
}

// Function to populate dropdowns
function populateDropdowns(characters) {
  const labelKeyDropdown = document.getElementById('labelKeyDropdown');
  const valueKeyDropdown = document.getElementById('valueKeyDropdown');
  const keys = Object.keys(characters[0]); // Assume all characters have the same keys

  keys.forEach(key => {
    const labelOption = document.createElement('option');
    labelOption.value = key;
    labelOption.textContent = key;
    labelKeyDropdown.appendChild(labelOption);

    const valueOption = document.createElement('option');
    valueOption.value = key;
    valueOption.textContent = key;
    valueKeyDropdown.appendChild(valueOption);
  });
}

document.addEventListener('DOMContentLoaded', function () {
  fetch('/data/allThePaths')
  .then(response => response.json())
  .then(data => {
    populateDropdowns(data.adventures)
    const container = document.getElementById("timelineigpaths");
    const items = new vis.DataSet(
      data.adventures.map((adventure, index) => {
        return {
          id: index + 1,
          content: adventure.adventure,
          start: adventure.igstart,
          end: adventure.igend,
          className: getColor(adventure.pathnr),
          customInfo: `${adventure.shortIntro}<hr>${adventure.adventureBackground}<hr>${adventure.otherBackground}`
        };
      })
    );
    const options = {};
    const timeline = new vis.Timeline(container, items, options);
    // Listen for select events
    timeline.on('select', function (properties) {
      const selectedItem = items.get(properties.items[0]); // Get the first selected item
      if (selectedItem) {
          // Update the info box with information from the selected item
          document.getElementById('infoBox').innerHTML = selectedItem.customInfo;
      } else {
          // If no item is selected, you can clear the info box or display a default message
          document.getElementById('infoBox').innerHTML = 'Click on an item to see more information here.';
      }
    });
  })
  .catch(error => {
      console.error('Error fetching adventure data:', error);
  });

  fetch('/data/allThePaths')
  .then(response => response.json())
  .then(data => {
    populateDropdowns(data.adventures)
    const container = document.getElementById("timelineirlpaths");
    const items = new vis.DataSet(
      data.adventures.map((adventure, index) => {
        return {
          id: index + 1,
          content: adventure.adventure,
          start: adventure.irlstart,
          end: adventure.irlend,
          className: getColor(adventure.pathnr),
          customInfo: `${adventure.shortIntro}<hr>${adventure.adventureBackground}<hr>${adventure.otherBackground}`
        };
      })
    );
    const options = {};
    const timeline = new vis.Timeline(container, items, options);
    // Listen for select events
    timeline.on('select', function (properties) {
      const selectedItem = items.get(properties.items[0]); // Get the first selected item
      if (selectedItem) {
          // Update the info box with information from the selected item
          document.getElementById('infoBox').innerHTML = selectedItem.customInfo;
      } else {
          // If no item is selected, you can clear the info box or display a default message
          document.getElementById('infoBox').innerHTML = 'Click on an item to see more information here.';
      }
    });
  })
  .catch(error => {
      console.error('Error fetching adventure data:', error);
  });

  fetch('/data/allTheData')
  .then(response => response.json())
  .then(data => {
    populateDropdowns(data.characters)
    const container = document.getElementById("timelineig");
    const items = new vis.DataSet(
      data.characters.map((character, index) => {
        return {
          id: index + 1,
          content: character.name,
          start: character.igstart,
          end: character.igend,
          className: getColor(character.player),
          customInfo: `${character.info}`
        };
      })
    );
    const options = {};
    const timeline = new vis.Timeline(container, items, options);
    // Listen for select events
    timeline.on('select', function (properties) {
      const selectedItem = items.get(properties.items[0]); // Get the first selected item
      if (selectedItem) {
          // Update the info box with information from the selected item
          document.getElementById('infoBox').innerHTML = selectedItem.customInfo;
      } else {
          // If no item is selected, you can clear the info box or display a default message
          document.getElementById('infoBox').innerHTML = 'Click on an item to see more information here.';
      }
    });
  })
  .catch(error => {
      console.error('Error fetching character data:', error);
  });

  fetch('/data/allTheData')
  .then(response => response.json())
  .then(data => {
    const container = document.getElementById("timelineirl");
    const items = new vis.DataSet(
      data.characters.map((character, index) => ({
        id: index + 1,
        content: character.name,
        start: character.irlstart,
        end: character.irlend,
        className: getColor(character.player),
        customInfo: `${character.info}`
      }))
    );
    const options = {};
    const timeline = new vis.Timeline(container, items, options);
    // Listen for select events
    timeline.on('select', function (properties) {
      const selectedItem = items.get(properties.items[0]); // Get the first selected item
      if (selectedItem) {
          // Update the info box with information from the selected item
          document.getElementById('infoBox').innerHTML = selectedItem.customInfo;
      } else {
          // If no item is selected, you can clear the info box or display a default message
          document.getElementById('infoBox').innerHTML = 'Click on an item to see more information here.';
      }
    });
  })
  .catch(error => {
      console.error('Error fetching character data:', error);
  });

});

document.addEventListener('DOMContentLoaded', function () {
  document.getElementById('generateChartButton').addEventListener('click', () => {
    const labelKey = document.getElementById('labelKeyDropdown').value;
    const valueKey = document.getElementById('valueKeyDropdown').value;

    // Hide all charts
    const charts = document.querySelectorAll('.chart');
    charts.forEach(chart => {
        chart.style.display = 'none';
    });
    const chartToShow = document.getElementById('myChart');
    chartToShow.style.display = 'block';

    createChartCount(dataUrl, 'myChart', labelKey + ' // ' + valueKey, labelKey, [valueKey]);
  });
});

document.addEventListener('DOMContentLoaded', function () {
  document.getElementById('plotList').addEventListener('click', function(e) {
    if (e.target && e.target.nodeName === "LI") {
      // Hide all charts
      const charts = document.querySelectorAll('.chart');
      charts.forEach(chart => {
          chart.style.display = 'none';
      });

      // Show the clicked chart
      const targetId = e.target.getAttribute('data-target');
      const chartToShow = document.getElementById(targetId);
      if (chartToShow) {
        console.log('Showing chart:', targetId);
        // run the function corresponding to the clicked chart
        if (targetId === 'charactercount') {
          createChartSum('/data/allTheData', 'charactercount', 'Deaths per Path (and Player)','path', ['player'], 'deaths');
        } else if (targetId === 'deathsByPath') {
          createChartSum('/data/allTheData', 'deathsByPath', 'Deaths per Player (and Path)','player', ['path'], 'deaths');
        } else if (targetId === 'classesByPlayer') {
          createChartCount('/data/allTheData', 'classesByPlayer', 'Classes per Player (excl multi-class)','player', ['class1']);
        } else if (targetId === 'classesByPlayerMulti') {
          createChartCount('/data/allTheData', 'classesByPlayerMulti', 'Classes per Player (incl multi-class)', 'player', ['class1','class2']);
        } else if (targetId === 'levelsLived') {
          createChartCount('/data/allTheData', 'levelsLived', 'Levels Lived (Player)','levelslived', ['player']);
        } else if (targetId === 'levelsLivedPath') {
          createChartCount('/data/allTheData', 'levelsLivedPath', 'Levels Lived (Path)','levelslived', ['path']);
        } else if (targetId === 'characterTotalScore') {
          createChartScore('/data/allTheData', 'characterTotalScore', 'Character Score');
        } else if (targetId === 'characterScore') {
          createChartAverageScore('/data/allTheData', 'characterScore','Character Survivability (levelslived/survival/crlvldiff/extralives)','name', 'lifescore');
        } else if (targetId === 'classScore') {
          createChartAverageScore('/data/allTheData', 'classScore','Class Survivability (how "easy" the class is, levelslived/survival/crlvldiff/extralives)','totalclass', 'classaverage');
        } else if (targetId === 'pathScore') {
          createChartAverageScore('/data/allTheData', 'pathScore', 'Path Survivability (how "easy" the path was, levelslived/survival/crlvldiff/extralives)', 'path', 'pathaverage');
        } else if (targetId === 'playerScore') {
          createChartAverageScore('/data/allTheData', 'playerScore', 'Player Survivability (levelslived/survival/crlvldiff/extralives)', 'player', 'playeraverage');
        } else if (targetId === 'killerCloud') {
          createWordCloud('/data/allTheData', 'killerCloud', aggregateKillerData);
        } else if (targetId === 'classCloud') {
          createWordCloud('/data/allTheData', 'classCloud', aggregateClassData);
        } else if (targetId === 'specializationCloud') {
          createWordCloud('/data/allTheData', 'specializationCloud', aggregateSpecializationData);
        } else if (targetId === 'raceCloud') {
          createWordCloud('/data/allTheData', 'raceCloud', aggregateRaceData);
        } else if (targetId === 'classWheelByPlayerExclude') {
          createWheel('/data/allTheData', 'classWheelByPlayerExclude', 'Class Wheel (excl multi-class)', 'exclude');
        } else if (targetId === 'classWheelByPlayerInclude') {
          createWheel('/data/allTheData', 'classWheelByPlayerInclude', 'Class Wheel (incl multi-class)', 'include');

        } else if (targetId === 'pathIgtime') {
          createChartSum('/data/allThePaths', 'pathIgtime', 'Adventures and In-game time','path', ['advnr'], 'igtime', 'adventures');
        } else if (targetId === 'pathIrltime') {
          createChartSum('/data/allThePaths', 'pathIrltime', 'Adventures and IRL time','path', ['advnr'], 'irltime', 'adventures');
        }
        chartToShow.style.display = 'block';
      }
    }
  });
});

function createChartSum(dataUrl, canvasId, chartTitle, labelKey, valueKeys, sumKey, dataKey = 'characters') {
  fetch(dataUrl)
    .then(response => response.json())
    .then(data => {
      const dataset = data[dataKey]; // Use dataKey to select the data set
      const ctx = document.getElementById(canvasId).getContext('2d');
      const labelValuePairs = {};
      const uniqueValues = new Set();

      // Initialize labelValuePairs with an object for each label
      dataset.forEach(item => {
        const label = item[labelKey];
        if (!labelValuePairs[label]) {
          labelValuePairs[label] = {};
        }
        valueKeys.forEach(valueKey => {
          const value = item[valueKey];
          if (value) {
            uniqueValues.add(value);
            if (!labelValuePairs[label][value]) {
              labelValuePairs[label][value] = 0;
            }
            const increment = sumKey && item[sumKey] ? parseFloat(item[sumKey]) : 1;
            labelValuePairs[label][value] += increment;
          }
        });
      });

      const labelsSorted = Object.keys(labelValuePairs).sort();
      const uniqueValuesSorted = Array.from(uniqueValues).sort();

      // Map each unique value to a dataset
      const datasets = uniqueValuesSorted.map((uniqueValue, index) => {
        const colorIndex = index % backgroundColors.length;
        return {
          label: uniqueValue,
          backgroundColor: backgroundColors[colorIndex],
          borderColor: borderColors[colorIndex],
          borderWidth: 1,
          data: labelsSorted.map(label => labelValuePairs[label][uniqueValue] || 0),
        };
      });

      if (chartInstances[canvasId]) {
        chartInstances[canvasId].destroy();
      }

      // Create the chart
      chartInstances[canvasId] = new Chart(ctx, {
        type: 'bar',
        data: {
          labels: labelsSorted,
          datasets: datasets,
        },
        options: {
          responsive: true,
          scales: {
            x: { stacked: true },
            y: { stacked: true, beginAtZero: true }
          },
          plugins: {
            title: {
              display: true,
              text: chartTitle,
            },
            legend: {
              display: true,
            }
          }
        }
      });
    });
}

function createChartCount(dataUrl, canvasId, chartTitle, labelKey, valueKeys, dataKey = 'characters') {
  fetch(dataUrl)
    .then(response => response.json())
    .then(data => {
      const dataset = data[dataKey]; // Use the dataKey to select the dataset
      const ctx = document.getElementById(canvasId).getContext('2d');
      const labelValuePairs = {};
      const uniqueValues = new Set();

      console.log(data);

      // Initialize labelValuePairs with an object for each label
      dataset.forEach(item => {
        const label = item[labelKey];
        if (!labelValuePairs[label]) {
          labelValuePairs[label] = {};
        }
        valueKeys.forEach(valueKey => {
          const value = item[valueKey];
          if (value) {
            uniqueValues.add(value);
            if (!labelValuePairs[label][value]) {
              labelValuePairs[label][value] = 0;
            }
            labelValuePairs[label][value]++;
          }
        });
      });

      const labelsSorted = Object.keys(labelValuePairs).sort();
      const uniqueValuesSorted = Array.from(uniqueValues).sort();

      const datasets = uniqueValuesSorted.map((uniqueValue, index) => {
        const colorIndex = index % backgroundColors.length;
        return {
          label: uniqueValue,
          backgroundColor: backgroundColors[colorIndex],
          borderColor: borderColors[colorIndex],
          borderWidth: 1,
          data: labelsSorted.map(label => labelValuePairs[label][uniqueValue] || 0),
        };
      });

      if (chartInstances[canvasId]) {
        chartInstances[canvasId].destroy();
      }

      chartInstances[canvasId] = new Chart(ctx, {
        type: 'bar',
        data: {
          labels: labelsSorted,
          datasets: datasets,
        },
        options: {
          responsive: true,
          scales: {
            x: { stacked: true },
            y: { stacked: true, beginAtZero: true }
          },
          plugins: {
            title: {
              display: true,
              text: chartTitle,
            },
            legend: {
              display: true,
            }
          }
        }
      });
    });
}


function createChartAverageScore(dataUrl, canvasId, chartTitle, groupKey, averageKey) {
  fetch(dataUrl)
    .then(response => response.json())
    .then(data => {
      const characters = data.characters;
      const ctx = document.getElementById(canvasId).getContext('2d');

      // Sort characters by the specified path average in descending order
      characters.sort((a, b) => Number(b[averageKey]) - Number(a[averageKey]));
      const uniqueCharacters = characters.filter((character, index, self) =>
        index === self.findIndex((c) => c[groupKey] === character[groupKey])
      );

      if (chartInstances[canvasId]) {
        chartInstances[canvasId].destroy();
      }

      // Create the chart
      chartInstances[canvasId] = new Chart(ctx, {
        type: 'bar', // Use 'horizontalBar' for horizontal bar chart
        data: {
          labels: uniqueCharacters.map(character => character[groupKey]), // Use path as labels for y-axis
          datasets: [{
            label: 'Average Score',
            data: uniqueCharacters.map(character => Number(character[averageKey])), // Use path average as bar lengths
            backgroundColor: 'rgba(54, 162, 235, 0.2)', // A single color for all bars
            borderColor: 'rgba(54, 162, 235, 1)', // A single border color for all bars
            borderWidth: 1,
          }]
        },
        options: {
          indexAxis: 'y', // Set 'y' to make the bar chart horizontal
          responsive: true,
          scales: {
            x: {
              beginAtZero: true,
              max: 100, // Explicitly set maximum to 100
            }
          },
          plugins: {
            legend: {
              display: false,
            },
            title: {
              display: true,
              text: chartTitle,
            }
          }
        }
      });
    });
}

function createChartScore(dataUrl, canvasId, chartTitle) {
  fetch(dataUrl)
    .then(response => response.json())
    .then(data => {
      const characters = data.characters;
      const ctx = document.getElementById(canvasId).getContext('2d');

      // Sort characters by characterscore in descending order
      characters.sort((a, b) => Number(b.characterscore) - Number(a.characterscore));

      // Create datasets for each score type with fixed colors
      const datasets = [
        {
          label: 'Path Score',
          data: characters.map(character => Number(character.pathscore)),
          backgroundColor: 'rgba(255, 99, 132, 0.2)',
          borderColor: 'rgba(255, 99, 132, 1)',
          borderWidth: 1
        },
        {
          label: 'Class Score',
          data: characters.map(character => Number(character.classscore)),
          backgroundColor: 'rgba(54, 162, 235, 0.2)',
          borderColor: 'rgba(54, 162, 235, 1)',
          borderWidth: 1
        },
        {
          label: 'Personal Score',
          data: characters.map(character => Number(character.lifescore)),
          backgroundColor: 'rgba(75, 192, 192, 0.2)',
          borderColor: 'rgba(75, 192, 192, 1)',
          borderWidth: 1
        }
      ];

      if (chartInstances[canvasId]) {
        chartInstances[canvasId].destroy(); // Destroy the previous instance if it exists
      }

      // Create the chart
      chartInstances[canvasId] = new Chart(ctx, {
        type: 'bar',
        data: {
          labels: characters.map(character => character.name),
          datasets: datasets
        },
        options: {
          indexAxis: 'y', // Horizontal bar chart
          responsive: true,
          scales: {
            x: {
              stacked: true, // Stacked bar chart
              beginAtZero: true
            }
          },
          plugins: {
            tooltip: {
              callbacks: {
                label: function(context) {
                  const character = characters[context.dataIndex];
                  const scoreType = context.dataset.label;
                  const scoreValue = context.parsed.x;
                  return `${character.name} - ${scoreType}: ${scoreValue}`;
                }
              }
            },
            title: {
              display: true,
              text: chartTitle
            },
            legend: {
              display: true // Set to true if you want to display the legend
            }
          }
        }
      });
    });
}

function createWheel(dataUrl, canvasIdPrefix, chartTitlePrefix, classcolumns) {
  fetch(dataUrl)
    .then(response => response.json())
    .then(data => {
      const categories = {};
      const allClasses = new Set();
      const shouldIgnoreClasses = document.getElementById('ignoreClassesCheckbox').checked;
      const includeAllClasses = document.getElementById('includeAllClassesCheckbox').checked; // New checkbox state

      data.characters.forEach(character => {
        const player = character.player;
        const class1 = character.class1;
        const class2 = character.class2;

        // Add to allClasses if not ignoring
        if (!shouldIgnoreClasses || !ignoredClasses.includes(class1)) {
          allClasses.add(class1);
        }
        if (class2 && (!shouldIgnoreClasses || !ignoredClasses.includes(class2))) {
          allClasses.add(class2);
        }

        // Initialize player in categories object if it doesn't exist
        if (!categories[player]) {
          categories[player] = {};
        }

        // Increment count for class1
        if (class1 && (!shouldIgnoreClasses || !ignoredClasses.includes(class1))) {
          categories[player][class1] = (categories[player][class1] || 0) + 1;
        }

        // Increment count for class2
        if (classcolumns == "include" && class2 && class2 !== class1 && (!shouldIgnoreClasses || !ignoredClasses.includes(class2))) {
          categories[player][class2] = (categories[player][class2] || 0) + 1;
        }
      });

      // Convert allClasses set to an array and sort it if including all classes
      const allClassesSorted = includeAllClasses ? Array.from(allClasses).sort() : null;

      const filteredAllClasses = [...allClasses].filter(className => !shouldIgnoreClasses || !ignoredClasses.includes(className));

      // Iterate through each player to create a radar chart
      Object.keys(categories).forEach((player, index) => {
        const canvasId = `${canvasIdPrefix}-${index}-${classcolumns}`;
        let canvasElement = document.getElementById(canvasId);
        if (!canvasElement) {
          const container = document.getElementById("chartsContainer");
          canvasElement = document.createElement('canvas');
          canvasElement.id = canvasId;
          container.appendChild(canvasElement);
        }

        const ctx = canvasElement.getContext('2d');

        // Determine classLabels based on includeAllClasses checkbox
        const classLabels = includeAllClasses ? filteredAllClasses : Object.keys(categories[player]).filter(className => categories[player].hasOwnProperty(className) && (!shouldIgnoreClasses || !ignoredClasses.includes(className))).sort();
        const classCounts = classLabels.map(className => categories[player][className] || 0);

        if (chartInstances[canvasId]) {
          chartInstances[canvasId].destroy();
        }

        chartInstances[canvasId] = new Chart(ctx, {
          type: 'radar',
          data: {
            labels: classLabels,
            datasets: [{
              label: `${chartTitlePrefix} - ${player}`,
              data: classCounts,
              fill: true,
              backgroundColor: "rgba(54, 162, 235, 0.2)",
              borderColor: "rgba(54, 162, 235, 1)",
              pointBackgroundColor: "rgba(54, 162, 235, 1)",
              pointBorderColor: "#fff",
              pointHoverBackgroundColor: "#fff",
              pointHoverBorderColor: "rgba(54, 162, 235, 1)"
            }]
          },
          options: {
            scales: {
              r: {
                angleLines: {
                  display: false
                },
                min: 0, // Explicitly set minimum to 0
                ticks: {
                  // Additional tick configuration if needed
                  stepSize: 1, // Sets the step size between ticks
                }
              }
            },
            // ... other options
          }
        });
      });
    });
}


// word clouds

function aggregateKillerData(data) {
  const counts = {};

  data.characters.forEach(character => {
      const { killer, killer_old } = character;
      if (killer) {
          counts[killer] = (counts[killer] || 0) + 1;
      }
      if (killer_old) {
          counts[killer_old] = (counts[killer_old] || 0) + 1;
      }
  });

  return counts;
}

function aggregateClassData(data) {
  const counts = {};

  data.characters.forEach(character => {
      const { class1, class2 } = character;
      if (class1) {
          counts[class1] = (counts[class1] || 0) + 1;
      }
      if (class2) {
          counts[class2] = (counts[class2] || 0) + 1;
      }
  });

  return counts;
}

function aggregateRaceData(data) {
  const counts = {};

  data.characters.forEach(character => {
    const { race } = character;
    if (race) {
      counts[race] = (counts[race] || 0) + 1;
    }
  });

  return counts;
}

function aggregateSpecializationData(data) {
  const counts = {};

  data.characters.forEach(character => {
    const { spec1, spec2 } = character;
    if (spec1 && spec1 !== "<low lvl>") {
        counts[spec1] = (counts[spec1] || 0) + 1;
    }
    if (spec2 && spec2 !== "<low lvl>") {
        counts[spec2] = (counts[spec2] || 0) + 1;
    }
});

  return counts;
}

function createWordCloud(dataUrl, cloudId, aggregateFunction) {
  fetch(dataUrl)
    .then(response => response.json())
    .then(data => {
      const dataCounts = aggregateFunction(data);
      const wordEntries = Object.entries(dataCounts).map(([text, size]) => ({ text, size }));

      const layout = d3.layout.cloud()
        .size([800, 600])
        .words(wordEntries.map(d => ({text: d.text, size: d.size * 10}))) // Adjust size scaling as needed
        .padding(5)
        .rotate(() => ~~(Math.random() * 2) * 90)
        .font("Impact")
        .fontSize(d => d.size)
        .on("end", draw);

      layout.start();

      function draw(words) {
        d3.select(`#${cloudId}`).append("svg")
          .attr("width", layout.size()[0])
          .attr("height", layout.size()[1])
          .append("g")
          .attr("transform", "translate(" + layout.size()[0] / 2 + "," + layout.size()[1] / 2 + ")")
          .selectAll("text")
          .data(words)
          .enter().append("text")
          .style("font-size", d => d.size + "px")
          .style("font-family", "Impact")
          .attr("text-anchor", "middle")
          .attr("transform", d => "translate(" + [d.x, d.y] + ")rotate(" + d.rotate + ")")
          .text(d => d.text);
      }
    });
}
