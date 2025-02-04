document.addEventListener("DOMContentLoaded", () => {
    const yamlOutput = document.getElementById("yamlOutput");
  
    // Add dynamic fields with a remove button
    const addDynamicField = (containerId, fieldHtml) => {
      const container = document.getElementById(containerId);
      const div = document.createElement("div");
      div.classList.add("dynamic-field");
      div.innerHTML = fieldHtml + '<button type="button" class="remove-btn">X</button>';
      container.appendChild(div);
  
      // Add event listener to the remove button
      div.querySelector(".remove-btn").addEventListener("click", () => {
        div.remove();
      });
    };
  
    // Dynamic field buttons
    document.getElementById("addArgumentButton").addEventListener("click", () => {
      addDynamicField("dataTypesContainer", `
        <label>Field:</label>
        <input type="text" class="argField" placeholder="e.g., argc">
        <label>Type:</label>
        <input type="text" class="argType" placeholder="e.g., integer">
        <label>Description:</label>
        <input type="text" class="argDescription" placeholder="e.g., The number of arguments">
      `);
    });
  
    document.getElementById("addSettingButton").addEventListener("click", () => {
      addDynamicField("settingsContainer", `
        <label>Field Name:</label>
        <input type="text" class="settingsField" placeholder="e.g., count">
        <label>Type:</label>
        <input type="text" class="settingsType" placeholder="e.g., unsigned integer">
        <label>Description:</label>
        <input type="text" class="settingsDescription" placeholder="e.g., Number of times to display the message">
      `);
    });
  
    document.getElementById("addFunctionButton").addEventListener("click", () => {
      addDynamicField("functionsContainer", `
        <label>Function Name:</label>
        <input type="text" class="functionName" placeholder="e.g., parse_arguments">
        <label>Description:</label>
        <input type="text" class="functionDescription" placeholder="e.g., Parse the command-line arguments">
        <label>Parameters:</label>
        <input type="text" class="functionParam" placeholder="e.g., ctx (Context)">
        <label>Return Value:</label>
        <input type="text" class="functionReturnValue" placeholder="e.g., HANDLE_ARGS">
      `);
    });
  
    document.getElementById("addStateButton").addEventListener("click", () => {
      addDynamicField("statesContainer", `
        <label>State Name:</label>
        <input type="text" class="stateName" placeholder="e.g., PARSE_ARGS">
        <label>Description:</label>
        <input type="text" class="stateDescription" placeholder="e.g., Parse command-line arguments">
      `);
    });
  
    document.getElementById("addStateTransitionButton").addEventListener("click", () => {
      addDynamicField("stateTableContainer", `
        <label>From State:</label>
        <input type="text" class="fromState" placeholder="e.g., START">
        <label>To State:</label>
        <input type="text" class="toState" placeholder="e.g., PARSE_ARGS">
        <label>Function:</label>
        <input type="text" class="transitionFunction" placeholder="e.g., parse_arguments">
      `);
    });
  
    document.getElementById("addPseudocodeButton").addEventListener("click", () => {
      addDynamicField("pseudocodeContainer", `
        <label>Function Name:</label>
        <input type="text" class="pseudocodeFunctionName" placeholder="e.g., parse_arguments">
        <label>Pseudocode:</label>
        <textarea class="pseudocodeText" placeholder="Write pseudocode here..."></textarea>
      `);
    });
  
    // Generate YAML
    document.getElementById("generateYamlButton").addEventListener("click", () => {
      try {
        // Get Purpose
        const purpose = document.getElementById("purpose").value.trim();
  
        // Get Data Types (Arguments)
        const argumentsData = Array.from(document.querySelectorAll("#dataTypesContainer .dataType")).map((arg) => ({
          field: arg.querySelector(".argField").value.trim(),
          type: arg.querySelector(".argType").value.trim(),
          description: arg.querySelector(".argDescription").value.trim(),
        }));
  
        // Get Settings
        const settings = Array.from(document.querySelectorAll("#settingsContainer .setting")).map((setting) => ({
          field: setting.querySelector(".settingsField").value.trim(),
          type: setting.querySelector(".settingsType").value.trim(),
          description: setting.querySelector(".settingsDescription").value.trim(),
        }));
  
        // Get Functions
        const functions = Array.from(document.querySelectorAll("#functionsContainer .function")).map((func) => ({
          name: func.querySelector(".functionName").value.trim(),
          description: func.querySelector(".functionDescription").value.trim(),
          parameters: func.querySelector(".functionParam").value.trim(),
          return: func.querySelector(".functionReturnValue").value.trim(),
        }));
  
        // Get States
        const states = Array.from(document.querySelectorAll("#statesContainer .state")).map((state) => ({
          name: state.querySelector(".stateName").value.trim(),
          description: state.querySelector(".stateDescription").value.trim(),
        }));
  
        // Get State Table
        const stateTable = Array.from(document.querySelectorAll("#stateTableContainer .stateTransition")).map((transition) => ({
          from: transition.querySelector(".fromState").value.trim(),
          to: transition.querySelector(".toState").value.trim(),
          function: transition.querySelector(".transitionFunction").value.trim(),
        }));
  
        // Get Pseudocode
        const pseudocode = Array.from(document.querySelectorAll("#pseudocodeContainer .pseudocode")).map((pseudo) => ({
          function: pseudo.querySelector(".pseudocodeFunctionName").value.trim(),
          code: pseudo.querySelector(".pseudocodeText").value.trim(),
        }));
  
        // Compile YAML Data
        const yamlData = {
          purpose,
          data_types: { arguments: argumentsData },
          settings,
          functions,
          states,
          state_table: stateTable,
          pseudocode,
        };
  
        // Generate YAML String
        const yamlString = jsyaml.dump(yamlData);
        yamlOutput.textContent = yamlString;
      } catch (error) {
        console.error("Error generating YAML:", error);
        alert("An error occurred while generating the YAML. Check the console for details.");
      }
    });
  
    // Generate DOT syntax for the state diagram
    function generateDotSyntax(stateTable) {
      let dot = "digraph G {\n";
      stateTable.forEach((transition) => {
        dot += `  "${transition.from}" -> "${transition.to}" [label="${transition.function}"];\n`;
      });
      dot += "}";
      return dot;
    }
  
    // Render the diagram as an SVG
    function renderDiagram(dotSyntax) {
      const viz = new Viz();
      return viz.renderSVGElement(dotSyntax);
    }
  
    // Convert SVG to PNG
    async function svgToPng(svgElement) {
      const svgData = new XMLSerializer().serializeToString(svgElement);
      const canvas = document.createElement("canvas");
      const context = canvas.getContext("2d");
  
      const svgBlob = new Blob([svgData], { type: "image/svg+xml;charset=utf-8" });
      const url = URL.createObjectURL(svgBlob);
  
      const img = new Image();
      img.src = url;
  
      return new Promise((resolve, reject) => {
        img.onload = () => {
          canvas.width = img.width;
          canvas.height = img.height;
          context.drawImage(img, 0, 0);
          URL.revokeObjectURL(url);
  
          canvas.toBlob(
            (blob) => {
              const reader = new FileReader();
              reader.onloadend = () => {
                const base64Data = reader.result.split(",")[1];
                resolve(base64ToUint8Array(base64Data));
              };
              reader.readAsDataURL(blob);
            },
            "image/png"
          );
        };
  
        img.onerror = (err) => reject(err);
      });
    }
  
    // Helper: Convert Base64 to Uint8Array
    function base64ToUint8Array(base64) {
      const binaryString = atob(base64);
      const len = binaryString.length;
      const bytes = new Uint8Array(len);
      for (let i = 0; i < len; i++) {
        bytes[i] = binaryString.charCodeAt(i);
      }
      return bytes;
    }
  
    // Download Word document
    document.getElementById("downloadWordButton").addEventListener("click", async () => {
        try {
          const yamlString = yamlOutput.textContent;
          if (!yamlString) {
            alert("Please generate the YAML first!");
            return;
          }
      
          const yamlData = jsyaml.load(yamlString);
      
          let pngData = null;
          if (yamlData.state_table && yamlData.state_table.length > 0) {
            const dotSyntax = generateDotSyntax(yamlData.state_table);
      
            try {
              const svgElement = await renderDiagram(dotSyntax);
              pngData = await svgToPng(svgElement);
            } catch (error) {
              console.error("Error rendering or converting diagram:", error);
              alert("Diagram generation failed, but the Word document will still be generated.");
            }
          }
      
          const { Document, Packer, Paragraph, TextRun, ImageRun } = window.docx;
      
          const doc = new Document({
            sections: [
              {
                children: [
                  // Document Title
                  new Paragraph({ text: "Assignment Design Document", heading: "Title" }),
                  new Paragraph({ text: "COMP 1234", heading: "Heading1" }),
                  new Paragraph({ text: "Assignment 1", heading: "Heading1" }),
                  new Paragraph({ text: "Design", heading: "Heading2" }),
                  new Paragraph({ text: "Pat Doe\nA00000000\nFeb 30th, 2020", spacing: { after: 200 } }),
      
                  // Purpose Section
                  new Paragraph({ text: "Purpose", heading: "Heading1" }),
                  new Paragraph({ text: yamlData.purpose || "N/A" }),
      
                  // Data Types Section
                  new Paragraph({ text: "Data Types", heading: "Heading1" }),
                  new Paragraph({ text: "Arguments", heading: "Heading2" }),
                  ...(yamlData.data_types.arguments.map((arg) =>
                    new Paragraph({
                      children: [
                        new TextRun({ text: `Field: ${arg.field}`, bold: true }),
                        new TextRun(`\nType: ${arg.type}`),
                        new TextRun(`\nDescription: ${arg.description}`),
                      ],
                    })
                  )),
      
                  // Settings Section
                  new Paragraph({ text: "Settings", heading: "Heading1" }),
                  ...(yamlData.settings.map((setting) =>
                    new Paragraph({
                      children: [
                        new TextRun({ text: `Field: ${setting.field}`, bold: true }),
                        new TextRun(`\nType: ${setting.type}`),
                        new TextRun(`\nDescription: ${setting.description}`),
                      ],
                    })
                  )),
      
                  // Functions Section
                  new Paragraph({ text: "Functions", heading: "Heading1" }),
                  ...(yamlData.functions.map((func) =>
                    new Paragraph({
                      children: [
                        new TextRun({ text: `Function Name: ${func.name}`, bold: true }),
                        new TextRun(`\nDescription: ${func.description}`),
                        new TextRun(`\nParameters: ${func.parameters}`),
                        new TextRun(`\nReturn: ${func.return}`),
                      ],
                    })
                  )),
      
                  // States Section
                  new Paragraph({ text: "States", heading: "Heading1" }),
                  ...(yamlData.states.map((state) =>
                    new Paragraph({
                      children: [
                        new TextRun({ text: `State: ${state.name}`, bold: true }),
                        new TextRun(`\nDescription: ${state.description}`),
                      ],
                    })
                  )),
      
                  // State Table Section
                  new Paragraph({ text: "State Table", heading: "Heading1" }),
                  ...(yamlData.state_table.map((transition) =>
                    new Paragraph({
                      children: [
                        new TextRun({ text: `From: ${transition.from}`, bold: true }),
                        new TextRun(`\nTo: ${transition.to}`),
                        new TextRun(`\nFunction: ${transition.function}`),
                      ],
                    })
                  )),
      
                  // State Diagram Section
                  ...(pngData
                    ? [
                        new Paragraph({ text: "State Transition Diagram", heading: "Heading1" }),
                        new Paragraph({
                          children: [
                            new ImageRun({
                              data: pngData,
                              transformation: { width: 600, height: 300 },
                            }),
                          ],
                        }),
                      ]
                    : []),
      
                  // Pseudocode Section
                  new Paragraph({ text: "Pseudocode", heading: "Heading1" }),
                  ...(yamlData.pseudocode.map((pseudo) =>
                    new Paragraph({
                      children: [
                        new TextRun({ text: `Function: ${pseudo.function}`, bold: true }),
                        new TextRun(`\nCode: ${pseudo.code}`),
                      ],
                    })
                  )),
                ],
              },
            ],
          });
      
          const blob = await Packer.toBlob(doc);
      
          const link = document.createElement("a");
          link.href = URL.createObjectURL(blob);
          link.download = "assignment_design.docx";
          link.click();
        } catch (error) {
          console.error("Error generating Word document:", error);
          alert("An error occurred while generating the Word document. Check the console for details.");
        }
      });
      
  });
  