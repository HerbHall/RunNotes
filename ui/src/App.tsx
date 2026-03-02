import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";

export default function App() {
  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom>
        RunNotes
      </Typography>
      <Typography variant="body1" color="text.secondary">
        Attach notes and annotations to your Docker containers.
      </Typography>
    </Box>
  );
}
