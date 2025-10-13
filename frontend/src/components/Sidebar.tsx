"use client";


import Box from '@mui/material/Box';
import Drawer from '@mui/material/Drawer';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemText from '@mui/material/ListItemText';
import { createContext, useContext, useState } from 'react';

type SidebarContextValue = {
  isOpen: boolean
  open: () => void
  close: () => void
  toggle: () => void
}


const SidebarContext = createContext<SidebarContextValue>({isOpen: false, open: () => {}, close: () => {}, toggle: () => {}})


export function useSidebar() {
  return useContext(SidebarContext);
}

export function SidebarProvider({children} : Readonly<{children: React.ReactNode}>) {
  const [isOpen, setIsOpen] = useState(false)
  
  const open = () => {
    setIsOpen(true)
  }
  const close = () => {
    setIsOpen(false)
  }
  const toggle = () => {
    setIsOpen(o => !o)
  }
  return <SidebarContext.Provider value={{isOpen, open, close, toggle }}>{children}</SidebarContext.Provider>
}

export default function Sidebar() {
  const sidebar = useContext(SidebarContext);
  return (
      <Drawer open={sidebar.isOpen} onClose={sidebar.close}>
        <Box sx={{ width: 200}}>
          <List>
            <ListItem disablePadding>
              <ListItemButton>
                <ListItemText>Services</ListItemText>
              </ListItemButton>
            </ListItem>
            <ListItem disablePadding>
              <ListItemButton>
                <ListItemText>Jobs</ListItemText>
              </ListItemButton>
            </ListItem>
          </List>
        </Box>
       
      </Drawer>
  );
}