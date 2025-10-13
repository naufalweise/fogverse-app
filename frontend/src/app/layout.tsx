import * as React from 'react';
import { AppRouterCacheProvider } from '@mui/material-nextjs/v15-appRouter';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import InitColorSchemeScript from '@mui/material/InitColorSchemeScript';
import theme from '@/theme';
import AppBar from "@/components/AppBar"
import ModeSwitch from '@/components/ModeSwitch';
import Sidebar, { SidebarProvider } from '@/components/Sidebar';


export default function RootLayout(props: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body>
        <InitColorSchemeScript attribute="class" />
        <AppRouterCacheProvider options={{ enableCssLayer: true }}>
          <ThemeProvider theme={theme}>
            <SidebarProvider>
              <CssBaseline />
              <AppBar />
              <ModeSwitch />
              <Sidebar />
              {props.children}
            </SidebarProvider>
          </ThemeProvider>
        </AppRouterCacheProvider>
      </body>
    </html>
  );
}