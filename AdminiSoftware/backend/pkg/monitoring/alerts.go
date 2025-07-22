
package monitoring

import (
	"fmt"
	"time"
)

type AlertManager struct {
	rules   []AlertRule
	alerts  []Alert
	channel chan Alert
}

type AlertRule struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Metric      string                 `json:"metric"`
	Condition   string                 `json:"condition"`
	Threshold   float64                `json:"threshold"`
	Duration    time.Duration          `json:"duration"`
	Severity    string                 `json:"severity"`
	Enabled     bool                   `json:"enabled"`
	Actions     []AlertAction          `json:"actions"`
	Labels      map[string]string      `json:"labels"`
}

type Alert struct {
	ID          string                 `json:"id"`
	RuleID      string                 `json:"rule_id"`
	Name        string                 `json:"name"`
	Message     string                 `json:"message"`
	Severity    string                 `json:"severity"`
	Status      string                 `json:"status"`
	StartsAt    time.Time              `json:"starts_at"`
	EndsAt      *time.Time             `json:"ends_at,omitempty"`
	Labels      map[string]string      `json:"labels"`
	Annotations map[string]string      `json:"annotations"`
}

type AlertAction struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

func NewAlertManager() *AlertManager {
	return &AlertManager{
		rules:   []AlertRule{},
		alerts:  []Alert{},
		channel: make(chan Alert, 100),
	}
}

func (am *AlertManager) AddRule(rule AlertRule) {
	am.rules = append(am.rules, rule)
}

func (am *AlertManager) GetRules() []AlertRule {
	return am.rules
}

func (am *AlertManager) GetActiveAlerts() []Alert {
	var active []Alert
	for _, alert := range am.alerts {
		if alert.Status == "firing" {
			active = append(active, alert)
		}
	}
	return active
}

func (am *AlertManager) EvaluateRules(metrics map[string]float64) {
	for _, rule := range am.rules {
		if !rule.Enabled {
			continue
		}
		
		value, exists := metrics[rule.Metric]
		if !exists {
			continue
		}
		
		triggered := am.evaluateCondition(value, rule.Condition, rule.Threshold)
		
		if triggered {
			am.fireAlert(rule, value)
		} else {
			am.resolveAlert(rule.ID)
		}
	}
}

func (am *AlertManager) evaluateCondition(value float64, condition string, threshold float64) bool {
	switch condition {
	case "greater_than":
		return value > threshold
	case "less_than":
		return value < threshold
	case "equal":
		return value == threshold
	case "greater_equal":
		return value >= threshold
	case "less_equal":
		return value <= threshold
	default:
		return false
	}
}

func (am *AlertManager) fireAlert(rule AlertRule, value float64) {
	// Check if alert already exists
	for i, alert := range am.alerts {
		if alert.RuleID == rule.ID && alert.Status == "firing" {
			return // Alert already firing
		}
		if alert.RuleID == rule.ID && alert.Status == "resolved" {
			// Refire the alert
			am.alerts[i].Status = "firing"
			am.alerts[i].StartsAt = time.Now()
			am.alerts[i].EndsAt = nil
			return
		}
	}
	
	alert := Alert{
		ID:       generateAlertID(),
		RuleID:   rule.ID,
		Name:     rule.Name,
		Message:  fmt.Sprintf("%s is %s %.2f (threshold: %.2f)", rule.Metric, rule.Condition, value, rule.Threshold),
		Severity: rule.Severity,
		Status:   "firing",
		StartsAt: time.Now(),
		Labels:   rule.Labels,
		Annotations: map[string]string{
			"metric":    rule.Metric,
			"value":     fmt.Sprintf("%.2f", value),
			"threshold": fmt.Sprintf("%.2f", rule.Threshold),
		},
	}
	
	am.alerts = append(am.alerts, alert)
	am.channel <- alert
	
	// Execute alert actions
	for _, action := range rule.Actions {
		am.executeAction(action, alert)
	}
}

func (am *AlertManager) resolveAlert(ruleID string) {
	for i, alert := range am.alerts {
		if alert.RuleID == ruleID && alert.Status == "firing" {
			now := time.Now()
			am.alerts[i].Status = "resolved"
			am.alerts[i].EndsAt = &now
			break
		}
	}
}

func (am *AlertManager) executeAction(action AlertAction, alert Alert) {
	switch action.Type {
	case "email":
		am.sendEmailAlert(action.Config, alert)
	case "webhook":
		am.sendWebhookAlert(action.Config, alert)
	case "slack":
		am.sendSlackAlert(action.Config, alert)
	}
}

func (am *AlertManager) sendEmailAlert(config map[string]interface{}, alert Alert) {
	// Implementation for email alerts
	fmt.Printf("Sending email alert: %s\n", alert.Message)
}

func (am *AlertManager) sendWebhookAlert(config map[string]interface{}, alert Alert) {
	// Implementation for webhook alerts
	fmt.Printf("Sending webhook alert: %s\n", alert.Message)
}

func (am *AlertManager) sendSlackAlert(config map[string]interface{}, alert Alert) {
	// Implementation for Slack alerts
	fmt.Printf("Sending Slack alert: %s\n", alert.Message)
}

func generateAlertID() string {
	return fmt.Sprintf("alert_%d", time.Now().UnixNano())
}

func (am *AlertManager) GetDefaultRules() []AlertRule {
	return []AlertRule{
		{
			ID:        "cpu_high",
			Name:      "High CPU Usage",
			Metric:    "cpu_usage",
			Condition: "greater_than",
			Threshold: 80.0,
			Duration:  5 * time.Minute,
			Severity:  "warning",
			Enabled:   true,
			Labels:    map[string]string{"type": "system"},
		},
		{
			ID:        "memory_high",
			Name:      "High Memory Usage",
			Metric:    "memory_usage",
			Condition: "greater_than",
			Threshold: 90.0,
			Duration:  5 * time.Minute,
			Severity:  "critical",
			Enabled:   true,
			Labels:    map[string]string{"type": "system"},
		},
		{
			ID:        "disk_full",
			Name:      "Disk Space Low",
			Metric:    "disk_usage",
			Condition: "greater_than",
			Threshold: 85.0,
			Duration:  10 * time.Minute,
			Severity:  "warning",
			Enabled:   true,
			Labels:    map[string]string{"type": "storage"},
		},
	}
}
