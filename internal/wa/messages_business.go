package wa

import (
	"strings"

	waProto "go.mau.fi/whatsmeow/binary/proto"
)

func extractBusinessText(m *waProto.Message, pm *ParsedMessage) {
	if tmpl := m.GetTemplateMessage(); tmpl != nil && pm.Text == "" {
		if hydrated := hydratedTemplate(tmpl); hydrated != nil {
			var parts []string
			if t := strings.TrimSpace(hydrated.GetHydratedTitleText()); t != "" {
				parts = append(parts, t)
			}
			if b := strings.TrimSpace(hydrated.GetHydratedContentText()); b != "" {
				parts = append(parts, b)
			}
			if f := strings.TrimSpace(hydrated.GetHydratedFooterText()); f != "" {
				parts = append(parts, "["+f+"]")
			}
			pm.Text = strings.Join(parts, "\n")
		} else if im := tmpl.GetInteractiveMessageTemplate(); im != nil {
			pm.Text = interactiveText(im)
		}
	}

	if btn := m.GetButtonsMessage(); btn != nil && pm.Text == "" {
		var parts []string
		if t := strings.TrimSpace(btn.GetText()); t != "" {
			parts = append(parts, t)
		}
		if b := strings.TrimSpace(btn.GetContentText()); b != "" {
			parts = append(parts, b)
		}
		if f := strings.TrimSpace(btn.GetFooterText()); f != "" {
			parts = append(parts, "["+f+"]")
		}
		pm.Text = strings.Join(parts, "\n")
	}

	if resp := m.GetButtonsResponseMessage(); resp != nil && pm.Text == "" {
		pm.Text = resp.GetSelectedDisplayText()
	}

	if im := m.GetInteractiveMessage(); im != nil && pm.Text == "" {
		pm.Text = interactiveText(im)
	}

	if resp := m.GetInteractiveResponseMessage(); resp != nil && pm.Text == "" {
		if body := resp.GetBody(); body != nil {
			pm.Text = strings.TrimSpace(body.GetText())
		}
	}

	if list := m.GetListMessage(); list != nil && pm.Text == "" {
		var parts []string
		if t := strings.TrimSpace(list.GetTitle()); t != "" {
			parts = append(parts, t)
		}
		if d := strings.TrimSpace(list.GetDescription()); d != "" {
			parts = append(parts, d)
		}
		pm.Text = strings.Join(parts, "\n")
	}

	if lr := m.GetListResponseMessage(); lr != nil && pm.Text == "" {
		pm.Text = strings.TrimSpace(lr.GetTitle())
		if pm.Text == "" {
			if sel := lr.GetSingleSelectReply(); sel != nil {
				pm.Text = sel.GetSelectedRowID()
			}
		}
	}

	if tbr := m.GetTemplateButtonReplyMessage(); tbr != nil && pm.Text == "" {
		pm.Text = tbr.GetSelectedDisplayText()
	}
}

func hydratedTemplate(tmpl *waProto.TemplateMessage) *waProto.TemplateMessage_HydratedFourRowTemplate {
	if h := tmpl.GetHydratedFourRowTemplate(); h != nil {
		return h
	}
	return tmpl.GetHydratedTemplate()
}

func interactiveText(im *waProto.InteractiveMessage) string {
	var parts []string
	if h := im.GetHeader(); h != nil {
		if t := strings.TrimSpace(h.GetTitle()); t != "" {
			parts = append(parts, t)
		}
	}
	if b := im.GetBody(); b != nil {
		if t := strings.TrimSpace(b.GetText()); t != "" {
			parts = append(parts, t)
		}
	}
	if f := im.GetFooter(); f != nil {
		if t := strings.TrimSpace(f.GetText()); t != "" {
			parts = append(parts, "["+t+"]")
		}
	}
	return strings.Join(parts, "\n")
}
