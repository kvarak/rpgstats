

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

document.addEventListener('DOMContentLoaded', function () {
  fetch('/data/allTheData')
  .then(response => response.json())
  .then(data => {
    const container = document.getElementById("timelineig");
    const items = new vis.DataSet(
      data.characters.map((character, index) => {
        return {
          id: index + 1,
          content: character.text,
          start: character.irlstart,
          end: character.irlend,
          className: getColor(character.category),
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
        content: character.text,
        start: character.irlstart,
        end: character.irlend,
        className: getColor(character.category),
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

// document.addEventListener('DOMContentLoaded', function () {
//   // Function to dynamically create checkbox HTML for each category
//   function createCategoryCheckboxes(categories) {
//     const container = document.querySelector('.checkbox-container');
//     // Clear existing checkboxes (if any)
//     container.innerHTML = '';

//     categories.forEach(category => {
//         const label = document.createElement('label');
//         label.innerHTML = `
//             <input type="checkbox" class="category-checkbox" value="${category}" checked> ${category}
//         `;
//         container.appendChild(label);
//     });
//   }

//   // Fetch the categories from the data
//   fetch('/data/charactercount')
//   .then(response => response.json())
//   .then(data => {
//     // Assuming data is an object where keys are category names
//     const categories = Object.keys(data);
//     createCategoryCheckboxes(categories);

//     // After creating checkboxes, you can also initialize the chart
//     // updateChart(...) or any other functionality you need
//   });
// });

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
          chartToShow.style.display = 'block';
      }
    }
  });
});

function createChart(dataUrl, canvasId, chartTitle) {
  fetch(dataUrl)
  .then(response => response.json())
  .then(data => {
    const ctx = document.getElementById(canvasId).getContext('2d');
    const allPaths = new Set();
    const names = Object.keys(data);

    // Aggregate all unique paths across all names
    names.forEach(name => {
      Object.keys(data[name]).forEach(path => {
        allPaths.add(path);
      });
    });

    const pathLabels = Array.from(allPaths).sort();
    const datasets = names.map((name, index) => {
      const colorIndex = index % backgroundColors.length; // Ensure this variable is accessible
      const dataForName = pathLabels.map(path => data[name][path] || 0);

      return {
        label: name,
        data: dataForName,
        backgroundColor: backgroundColors[colorIndex],
        borderColor: borderColors[colorIndex],
        borderWidth: 1
      };
    });

    new Chart(ctx, {
      type: 'bar',
      data: {
        labels: pathLabels,
        datasets: datasets
      },
      options: {
        responsive: false,
        scales: {
          x: { stacked: true },
          y: { stacked: true, beginAtZero: true }
        },
        plugins: {
          title: {
            display: true,
            text: chartTitle,
            font: { size: 20 },
            padding: { top: 10, bottom: 30 }
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
      // Initialize a categories object to hold the counts of each class1 and class2
      const categories = {};

      // Iterate through each character to populate the categories object
      data.characters.forEach(character => {
        const category = character.category;
        const class1 = character.class1;
        const class2 = character.class2;

        // Initialize category in categories object if it doesn't exist
        if (!categories[category]) {
          categories[category] = {};
        }

        // Increment count for class1
        if (class1) {
          categories[category][class1] = (categories[category][class1] || 0) + 1;
        }

        if (classcolumns == "include") {
          // Increment count for class2 if it's different from class1
          if (class2 && class2 !== class1) {
            categories[category][class2] = (categories[category][class2] || 0) + 1;
          }
        }
      });

      // Iterate through each category to create a radar chart
      Object.keys(categories).forEach((category, index) => {
        // Prepare canvas
        const canvasId = `${canvasIdPrefix}-${index}-${classcolumns}`;
        let canvasElement = document.getElementById(canvasId);
        if (!canvasElement) {
          // If canvas does not exist, create it and append to a container
          const container = document.getElementById("chartsContainer"); // Ensure this is your container's ID
          canvasElement = document.createElement('canvas');
          canvasElement.id = canvasId;
          container.appendChild(canvasElement);
        }

        const ctx = canvasElement.getContext('2d');

        // Prepare data for the radar chart
        const classLabels = Object.keys(categories[category]);
        const classCounts = classLabels.map(className => categories[category][className]);

        // Generate the radar chart
        new Chart(ctx, {
          type: 'radar',
          data: {
            labels: classLabels,
            datasets: [{
              label: `${chartTitlePrefix} - ${category}`,
              data: classCounts,
              fill: true,
              backgroundColor: "rgba(54, 162, 235, 0.2)", // Example background color
              borderColor: "rgba(54, 162, 235, 1)", // Example border color
              pointBackgroundColor: "rgba(54, 162, 235, 1)", // Example point color
              pointBorderColor: "#fff",
              pointHoverBackgroundColor: "#fff",
              pointHoverBorderColor: "rgba(54, 162, 235, 1)"
            }]
          },
          options: {
            elements: {
              line: {
                borderWidth: 3
              }
            },
            plugins: {
              title: {
                display: true,
                text: `${chartTitlePrefix} - ${category}`,
                font: { size: 20 },
                padding: { top: 10, bottom: 30 }
              },
              legend: {
                display: false
              }
            },
            scales: {
              r: {
                angleLines: {
                  display: false
                },
                suggestedMin: 0
              }
            }
          }
        });
      });
    });
}


// Killer word cloud

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
