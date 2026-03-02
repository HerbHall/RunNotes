import Box from "@mui/material/Box";
import ReactMarkdown from "react-markdown";

interface MarkdownPreviewProps {
  content: string;
}

export function MarkdownPreview({ content }: MarkdownPreviewProps) {
  return (
    <Box
      sx={{
        "& h1": { fontSize: "1.5rem", fontWeight: 600, mt: 2, mb: 1 },
        "& h2": { fontSize: "1.25rem", fontWeight: 600, mt: 2, mb: 1 },
        "& h3": { fontSize: "1.1rem", fontWeight: 600, mt: 1.5, mb: 0.5 },
        "& p": { mb: 1 },
        "& ul, & ol": { pl: 3, mb: 1 },
        "& code": {
          fontFamily: "monospace",
          backgroundColor: "action.hover",
          px: 0.5,
          borderRadius: 0.5,
          fontSize: "0.875em",
        },
        "& pre": {
          backgroundColor: "action.hover",
          p: 1.5,
          borderRadius: 1,
          overflow: "auto",
          mb: 1,
          "& code": { backgroundColor: "transparent", px: 0 },
        },
        "& blockquote": {
          borderLeft: 3,
          borderColor: "divider",
          pl: 2,
          ml: 0,
          color: "text.secondary",
        },
        "& a": { color: "primary.main" },
        "& hr": { borderColor: "divider", my: 2 },
      }}
    >
      <ReactMarkdown>{content}</ReactMarkdown>
    </Box>
  );
}
