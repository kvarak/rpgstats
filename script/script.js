
document.addEventListener('DOMContentLoaded', function () {
  // Function to dynamically create checkbox HTML for each category
  function createCategoryCheckboxes(categories) {
    const container = document.querySelector('.checkbox-container');
    // Clear existing checkboxes (if any)
    container.innerHTML = '';

    categories.forEach(category => {
        const label = document.createElement('label');
        label.innerHTML = `
            <input type="checkbox" class="category-checkbox" value="${category}" checked> ${category}
        `;
        container.appendChild(label);
    });
  }

  // Fetch the categories from the data
  fetch('/data/charactercount')
  .then(response => response.json())
  .then(data => {
    // Assuming data is an object where keys are category names
    const categories = Object.keys(data);
    createCategoryCheckboxes(categories);

    // After creating checkboxes, you can also initialize the chart
    // updateChart(...) or any other functionality you need
  });
});

// document.addEventListener('DOMContentLoaded', function () {
//   document.getElementById('plotList').addEventListener('click', function(e) {
//     if (e.target && e.target.nodeName == "LI") {
//       const targetId = e.target.getAttribute('data-target');
//       const chart = document.getElementById(targetId);
//       if (chart) {
//           chart.style.display = chart.style.display === 'none' ? 'block' : 'none';
//       }
//     }
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
