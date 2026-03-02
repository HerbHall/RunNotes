import Badge from "@mui/material/Badge";
import Box from "@mui/material/Box";
import CircularProgress from "@mui/material/CircularProgress";
import List from "@mui/material/List";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import Typography from "@mui/material/Typography";
import NoteIcon from "@mui/icons-material/Description";
import type { ContainerInfo } from "../types";

interface ContainerListProps {
  containers: ContainerInfo[];
  selectedName: string | null;
  onSelect: (name: string) => void;
  noteCounts: Map<string, number>;
  loading: boolean;
}

export function ContainerList({
  containers,
  selectedName,
  onSelect,
  noteCounts,
  loading,
}: ContainerListProps) {
  if (loading) {
    return (
      <Box
        sx={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          py: 4,
        }}
      >
        <CircularProgress size={32} />
      </Box>
    );
  }

  if (containers.length === 0) {
    return (
      <Box sx={{ px: 2, py: 4, textAlign: "center" }}>
        <Typography variant="body2" color="text.secondary">
          No containers found
        </Typography>
      </Box>
    );
  }

  return (
    <List disablePadding>
      {containers.map((c) => {
        const count = noteCounts.get(c.name) ?? 0;
        return (
          <ListItemButton
            key={c.id}
            selected={selectedName === c.name}
            onClick={() => onSelect(c.name)}
            sx={{ borderBottom: 1, borderColor: "divider" }}
          >
            <ListItemIcon sx={{ minWidth: 36 }}>
              <Box
                sx={{
                  width: 10,
                  height: 10,
                  borderRadius: "50%",
                  bgcolor: c.state === "running" ? "success.main" : "grey.500",
                }}
              />
            </ListItemIcon>
            <ListItemText
              primary={c.name}
              secondary={`${c.image} - ${c.status}`}
              primaryTypographyProps={{ noWrap: true }}
              secondaryTypographyProps={{ noWrap: true }}
            />
            {count > 0 && (
              count > 1 ? (
                <Badge badgeContent={count} color="primary" sx={{ ml: 1 }}>
                  <NoteIcon fontSize="small" color="primary" />
                </Badge>
              ) : (
                <NoteIcon fontSize="small" color="primary" sx={{ ml: 1 }} />
              )
            )}
          </ListItemButton>
        );
      })}
    </List>
  );
}
