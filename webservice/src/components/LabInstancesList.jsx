import React from 'react';

import Paper from '@material-ui/core/Paper';
import RunningLabList from './RunningLabList';

/* The style for the ListItem */
/**
 * Function to draw a list of running lab instances
 * @param props contains all the function to be associated with the components (buttons click, etc.)
 * @return The component to be drawn
 */
export default function LabInstancesList(props) {
  /* Parsing the instances array and draw for each one a list item with the right coloration, according to its status */
  const { runningLabs, stop, connect, isStudentView } = props;

  const runningLabNames = Array.from(runningLabs.keys());
  const runningLabList = runningLabNames.map(labName => ({
    ...runningLabs.get(labName),
    labName
  }));

  return (
    <Paper
      elevation={6}
      style={{
        flex: 1,
        minWidth: 575,
        maxWidth: 650,
        padding: 10,
        margin: 10,
        maxHeight: '70vh'
      }}
    >
      <RunningLabList
        labList={runningLabList}
        stop={stop}
        connect={connect}
        title="Running images"
        isStudentView={isStudentView}
      />
    </Paper>
  );
}
