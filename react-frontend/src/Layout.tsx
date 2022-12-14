import React from 'react';
import {
  AppBar,
  Box, Button, Container,
  Toolbar,
  Typography
} from "@mui/material";
import HomeView from "./components/Home";
import RecordView from "./components/Record";
import MockView from "./components/Mock";

const navItems = ['Home', 'Record', 'Mock'];

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;
  return (
    <div
      hidden={value !== index}
      id={`tabpanel-${index}`}
      aria-labelledby={`tab-${index}`}
      {...other}
    >
      { value === index && (children)}
    </div>
  );
}

export default function MainLayout() {
  const [tabValue, setTabValue] = React.useState(0);

  const handleChange = (newValue: number) => {
    setTabValue(newValue);
  };

  return (
    <Box sx={{ display: 'flex' }}>
      <AppBar component='nav'>
        <Toolbar>
          <Typography variant='h6' component='div' noWrap sx={{ mr: 2, display: 'flex' }}>
            Cocoon
          </Typography>
          <Box sx={{ flexGrow: 1, display: 'block' }}>

            {navItems.map((item, idx) => (
              <Button key={item} onClick={() => handleChange(idx) } sx={{ color: '#fff' }}>{item}</Button>
            ))}
          </Box>
        </Toolbar>
      </AppBar>
      <Container component='main' maxWidth={false} sx={{ p: 1, m: 0 }}>
        <Toolbar />
        <TabPanel value={tabValue} index={0}>
          <HomeView />
        </TabPanel>
        <TabPanel value={tabValue} index={1}>
          <RecordView />
        </TabPanel>
        <TabPanel index={tabValue} value={2}>
          <MockView />
        </TabPanel>
      </Container>
    </Box>
  )
}