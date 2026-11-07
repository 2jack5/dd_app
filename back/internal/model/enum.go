package model

// 状态字段使用独立类型；状态迁移与权限、期限校验放在 biz/data 的事务用例。

type UserStatus string

const (
	UserStatusActive        UserStatus = "ACTIVE"
	UserStatusBanned        UserStatus = "BANNED"
	UserStatusCancelPending UserStatus = "CANCEL_PENDING"
	UserStatusCancelled     UserStatus = "CANCELLED"
)

type AuthIdentityProvider string

const (
	AuthIdentityProviderWechatMini AuthIdentityProvider = "WECHAT_MINI"
	AuthIdentityProviderWechatApp  AuthIdentityProvider = "WECHAT_APP"
)

type AgreementAcceptanceAgreementCode string

const (
	AgreementAcceptanceAgreementCodeUser    AgreementAcceptanceAgreementCode = "USER"
	AgreementAcceptanceAgreementCodePrivacy AgreementAcceptanceAgreementCode = "PRIVACY"
	AgreementAcceptanceAgreementCodeCoach   AgreementAcceptanceAgreementCode = "COACH"
	AgreementAcceptanceAgreementCodeOrder   AgreementAcceptanceAgreementCode = "ORDER"
)

type MediaAssetPurpose string

const (
	MediaAssetPurposeAvatar         MediaAssetPurpose = "AVATAR"
	MediaAssetPurposeSkill          MediaAssetPurpose = "SKILL"
	MediaAssetPurposeReview         MediaAssetPurpose = "REVIEW"
	MediaAssetPurposeCase           MediaAssetPurpose = "CASE"
	MediaAssetPurposeBill           MediaAssetPurpose = "BILL"
	MediaAssetPurposeAgreement      MediaAssetPurpose = "AGREEMENT"
	MediaAssetPurposeChannelPayload MediaAssetPurpose = "CHANNEL_PAYLOAD"
	MediaAssetPurposeGameIcon       MediaAssetPurpose = "GAME_ICON"
	MediaAssetPurposeBanner         MediaAssetPurpose = "BANNER"
	MediaAssetPurposeActivityCover  MediaAssetPurpose = "ACTIVITY_COVER"
	MediaAssetPurposeServiceImage   MediaAssetPurpose = "SERVICE_IMAGE"
)

type MediaAssetAccessScope string

const (
	MediaAssetAccessScopePrivate MediaAssetAccessScope = "PRIVATE"
	MediaAssetAccessScopePublic  MediaAssetAccessScope = "PUBLIC"
)

type MediaAssetStatus string

const (
	MediaAssetStatusUploading MediaAssetStatus = "UPLOADING"
	MediaAssetStatusPending   MediaAssetStatus = "PENDING"
	MediaAssetStatusApproved  MediaAssetStatus = "APPROVED"
	MediaAssetStatusRejected  MediaAssetStatus = "REJECTED"
	MediaAssetStatusDeleted   MediaAssetStatus = "DELETED"
)

type NotificationBizType string

const (
	NotificationBizTypeApplication          NotificationBizType = "APPLICATION"
	NotificationBizTypeOrder                NotificationBizType = "ORDER"
	NotificationBizTypeRefund               NotificationBizType = "REFUND"
	NotificationBizTypeCase                 NotificationBizType = "CASE"
	NotificationBizTypeSettlement           NotificationBizType = "SETTLEMENT"
	NotificationBizTypeSettlementSupplement NotificationBizType = "SETTLEMENT_SUPPLEMENT"
)

type NotificationChannel string

const (
	NotificationChannelInApp           NotificationChannel = "IN_APP"
	NotificationChannelWechatSubscribe NotificationChannel = "WECHAT_SUBSCRIBE"
)

type NotificationStatus string

const (
	NotificationStatusPending NotificationStatus = "PENDING"
	NotificationStatusSent    NotificationStatus = "SENT"
	NotificationStatusFailed  NotificationStatus = "FAILED"
)

type CoachStatus string

const (
	CoachStatusApplying  CoachStatus = "APPLYING"
	CoachStatusActive    CoachStatus = "ACTIVE"
	CoachStatusSuspended CoachStatus = "SUSPENDED"
	CoachStatusRejected  CoachStatus = "REJECTED"
	CoachStatusClosed    CoachStatus = "CLOSED"
)

type CoachApplicationStatus string

const (
	CoachApplicationStatusDraft     CoachApplicationStatus = "DRAFT"
	CoachApplicationStatusSubmitted CoachApplicationStatus = "SUBMITTED"
	CoachApplicationStatusReviewing CoachApplicationStatus = "REVIEWING"
	CoachApplicationStatusApproved  CoachApplicationStatus = "APPROVED"
	CoachApplicationStatusRejected  CoachApplicationStatus = "REJECTED"
	CoachApplicationStatusWithdrawn CoachApplicationStatus = "WITHDRAWN"
)

type IdentityCheckStatus string

const (
	IdentityCheckStatusPending IdentityCheckStatus = "PENDING"
	IdentityCheckStatusPassed  IdentityCheckStatus = "PASSED"
	IdentityCheckStatusFailed  IdentityCheckStatus = "FAILED"
	IdentityCheckStatusExpired IdentityCheckStatus = "EXPIRED"
	IdentityCheckStatusRevoked IdentityCheckStatus = "REVOKED"
)

type IdentityCheckAgeCheckResult string

const (
	IdentityCheckAgeCheckResultUnknown    IdentityCheckAgeCheckResult = "UNKNOWN"
	IdentityCheckAgeCheckResultEligible   IdentityCheckAgeCheckResult = "ELIGIBLE"
	IdentityCheckAgeCheckResultIneligible IdentityCheckAgeCheckResult = "INELIGIBLE"
)

type CoachSkillStatus string

const (
	CoachSkillStatusDraft     CoachSkillStatus = "DRAFT"
	CoachSkillStatusSubmitted CoachSkillStatus = "SUBMITTED"
	CoachSkillStatusReviewing CoachSkillStatus = "REVIEWING"
	CoachSkillStatusApproved  CoachSkillStatus = "APPROVED"
	CoachSkillStatusRejected  CoachSkillStatus = "REJECTED"
	CoachSkillStatusExpired   CoachSkillStatus = "EXPIRED"
)

type GameStatus string

const (
	GameStatusEnabled  GameStatus = "ENABLED"
	GameStatusDisabled GameStatus = "DISABLED"
)

type GameRegionStatus string

const (
	GameRegionStatusEnabled  GameRegionStatus = "ENABLED"
	GameRegionStatusDisabled GameRegionStatus = "DISABLED"
)

type GameRankStatus string

const (
	GameRankStatusEnabled  GameRankStatus = "ENABLED"
	GameRankStatusDisabled GameRankStatus = "DISABLED"
)

type ServiceCategoryUnitType string

const (
	ServiceCategoryUnitTypeHour  ServiceCategoryUnitType = "HOUR"
	ServiceCategoryUnitTypeRound ServiceCategoryUnitType = "ROUND"
)

type ServiceCategoryStatus string

const (
	ServiceCategoryStatusEnabled  ServiceCategoryStatus = "ENABLED"
	ServiceCategoryStatusDisabled ServiceCategoryStatus = "DISABLED"
)

type CoachServiceUnitType string

const (
	CoachServiceUnitTypeHour  CoachServiceUnitType = "HOUR"
	CoachServiceUnitTypeRound CoachServiceUnitType = "ROUND"
)

type CoachServiceStatus string

const (
	CoachServiceStatusDraft    CoachServiceStatus = "DRAFT"
	CoachServiceStatusOnShelf  CoachServiceStatus = "ON_SHELF"
	CoachServiceStatusOffShelf CoachServiceStatus = "OFF_SHELF"
	CoachServiceStatusBlocked  CoachServiceStatus = "BLOCKED"
)

type CoachAvailabilityStatus string

const (
	CoachAvailabilityStatusOpen   CoachAvailabilityStatus = "OPEN"
	CoachAvailabilityStatusClosed CoachAvailabilityStatus = "CLOSED"
)

type BookingHoldStatus string

const (
	BookingHoldStatusHeld      BookingHoldStatus = "HELD"
	BookingHoldStatusConfirmed BookingHoldStatus = "CONFIRMED"
	BookingHoldStatusReleased  BookingHoldStatus = "RELEASED"
	BookingHoldStatusExpired   BookingHoldStatus = "EXPIRED"
)

type OrderQuoteStatus string

const (
	OrderQuoteStatusValid    OrderQuoteStatus = "VALID"
	OrderQuoteStatusConsumed OrderQuoteStatus = "CONSUMED"
	OrderQuoteStatusExpired  OrderQuoteStatus = "EXPIRED"
)

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusWaitAccept     OrderStatus = "WAIT_ACCEPT"
	OrderStatusWaitStart      OrderStatus = "WAIT_START"
	OrderStatusInService      OrderStatus = "IN_SERVICE"
	OrderStatusWaitConfirm    OrderStatus = "WAIT_CONFIRM"
	OrderStatusCompleted      OrderStatus = "COMPLETED"
	OrderStatusClosed         OrderStatus = "CLOSED"
)

type OrderDisputeStatus string

const (
	OrderDisputeStatusNone     OrderDisputeStatus = "NONE"
	OrderDisputeStatusOpen     OrderDisputeStatus = "OPEN"
	OrderDisputeStatusResolved OrderDisputeStatus = "RESOLVED"
)

type OrderEventActorType string

const (
	OrderEventActorTypeUser   OrderEventActorType = "USER"
	OrderEventActorTypeCoach  OrderEventActorType = "COACH"
	OrderEventActorTypeAdmin  OrderEventActorType = "ADMIN"
	OrderEventActorTypeSystem OrderEventActorType = "SYSTEM"
)

type PaymentStatus string

const (
	PaymentStatusUnpaid        PaymentStatus = "UNPAID"
	PaymentStatusPrepayCreated PaymentStatus = "PREPAY_CREATED"
	PaymentStatusPaid          PaymentStatus = "PAID"
	PaymentStatusClosing       PaymentStatus = "CLOSING"
	PaymentStatusClosed        PaymentStatus = "CLOSED"
	PaymentStatusFailed        PaymentStatus = "FAILED"
)

type PaymentEventEventType string

const (
	PaymentEventEventTypePayNotify PaymentEventEventType = "PAY_NOTIFY"
	PaymentEventEventTypeQuery     PaymentEventEventType = "QUERY"
	PaymentEventEventTypeClose     PaymentEventEventType = "CLOSE"
)

type PaymentEventStatus string

const (
	PaymentEventStatusReceived  PaymentEventStatus = "RECEIVED"
	PaymentEventStatusProcessed PaymentEventStatus = "PROCESSED"
	PaymentEventStatusRetry     PaymentEventStatus = "RETRY"
	PaymentEventStatusRejected  PaymentEventStatus = "REJECTED"
)

type RefundStatus string

const (
	RefundStatusApplied    RefundStatus = "APPLIED"
	RefundStatusReviewing  RefundStatus = "REVIEWING"
	RefundStatusApproved   RefundStatus = "APPROVED"
	RefundStatusSubmitted  RefundStatus = "SUBMITTED"
	RefundStatusProcessing RefundStatus = "PROCESSING"
	RefundStatusSuccess    RefundStatus = "SUCCESS"
	RefundStatusFailed     RefundStatus = "FAILED"
	RefundStatusRejected   RefundStatus = "REJECTED"
	RefundStatusCancelled  RefundStatus = "CANCELLED"
)

type RefundFundsStatus string

const (
	RefundFundsStatusNone     RefundFundsStatus = "NONE"
	RefundFundsStatusReserved RefundFundsStatus = "RESERVED"
	RefundFundsStatusConsumed RefundFundsStatus = "CONSUMED"
	RefundFundsStatusReleased RefundFundsStatus = "RELEASED"
)

type RefundEventActorType string

const (
	RefundEventActorTypeUser     RefundEventActorType = "USER"
	RefundEventActorTypeAdmin    RefundEventActorType = "ADMIN"
	RefundEventActorTypeProvider RefundEventActorType = "PROVIDER"
	RefundEventActorTypeSystem   RefundEventActorType = "SYSTEM"
)

type ReviewVisibility string

const (
	ReviewVisibilityPending   ReviewVisibility = "PENDING"
	ReviewVisibilityPublished ReviewVisibility = "PUBLISHED"
	ReviewVisibilityHidden    ReviewVisibility = "HIDDEN"
)

type ReviewContentKind string

const (
	ReviewContentKindPrimary  ReviewContentKind = "PRIMARY"
	ReviewContentKindFollowUp ReviewContentKind = "FOLLOW_UP"
)

type ReviewContentAuditStatus string

const (
	ReviewContentAuditStatusPending  ReviewContentAuditStatus = "PENDING"
	ReviewContentAuditStatusApproved ReviewContentAuditStatus = "APPROVED"
	ReviewContentAuditStatusRejected ReviewContentAuditStatus = "REJECTED"
)

type AfterSalesCaseKind string

const (
	AfterSalesCaseKindComplaint    AfterSalesCaseKind = "COMPLAINT"
	AfterSalesCaseKindSupport      AfterSalesCaseKind = "SUPPORT"
	AfterSalesCaseKindReviewAppeal AfterSalesCaseKind = "REVIEW_APPEAL"
)

type AfterSalesCaseStatus string

const (
	AfterSalesCaseStatusOpen          AfterSalesCaseStatus = "OPEN"
	AfterSalesCaseStatusAssigned      AfterSalesCaseStatus = "ASSIGNED"
	AfterSalesCaseStatusInvestigating AfterSalesCaseStatus = "INVESTIGATING"
	AfterSalesCaseStatusWaitUser      AfterSalesCaseStatus = "WAIT_USER"
	AfterSalesCaseStatusResolved      AfterSalesCaseStatus = "RESOLVED"
	AfterSalesCaseStatusClosed        AfterSalesCaseStatus = "CLOSED"
)

type CaseEventActorType string

const (
	CaseEventActorTypeUser   CaseEventActorType = "USER"
	CaseEventActorTypeAdmin  CaseEventActorType = "ADMIN"
	CaseEventActorTypeSystem CaseEventActorType = "SYSTEM"
)

type CaseEventEventType string

const (
	CaseEventEventTypeComment         CaseEventEventType = "COMMENT"
	CaseEventEventTypeAssign          CaseEventEventType = "ASSIGN"
	CaseEventEventTypeRequestEvidence CaseEventEventType = "REQUEST_EVIDENCE"
	CaseEventEventTypeResolve         CaseEventEventType = "RESOLVE"
	CaseEventEventTypeReopen          CaseEventEventType = "REOPEN"
)

type CaseEventVisibility string

const (
	CaseEventVisibilityPublic   CaseEventVisibility = "PUBLIC"
	CaseEventVisibilityInternal CaseEventVisibility = "INTERNAL"
)

type SettlementAccountStatus string

const (
	SettlementAccountStatusPending   SettlementAccountStatus = "PENDING"
	SettlementAccountStatusReviewing SettlementAccountStatus = "REVIEWING"
	SettlementAccountStatusActive    SettlementAccountStatus = "ACTIVE"
	SettlementAccountStatusRejected  SettlementAccountStatus = "REJECTED"
	SettlementAccountStatusDisabled  SettlementAccountStatus = "DISABLED"
)

type CommissionRuleScopeType string

const (
	CommissionRuleScopeTypeGlobal       CommissionRuleScopeType = "GLOBAL"
	CommissionRuleScopeTypeGameCategory CommissionRuleScopeType = "GAME_CATEGORY"
	CommissionRuleScopeTypeCoach        CommissionRuleScopeType = "COACH"
)

type CommissionRuleStatus string

const (
	CommissionRuleStatusDraft     CommissionRuleStatus = "DRAFT"
	CommissionRuleStatusPublished CommissionRuleStatus = "PUBLISHED"
	CommissionRuleStatusRetired   CommissionRuleStatus = "RETIRED"
)

type SettlementItemStatus string

const (
	SettlementItemStatusNotEligible SettlementItemStatus = "NOT_ELIGIBLE"
	SettlementItemStatusPending     SettlementItemStatus = "PENDING"
	SettlementItemStatusHolding     SettlementItemStatus = "HOLDING"
	SettlementItemStatusSubmitted   SettlementItemStatus = "SUBMITTED"
	SettlementItemStatusSettled     SettlementItemStatus = "SETTLED"
	SettlementItemStatusFailed      SettlementItemStatus = "FAILED"
	SettlementItemStatusNoPayable   SettlementItemStatus = "NO_PAYABLE"
)

type SettlementAdjustmentKind string

const (
	SettlementAdjustmentKindRefund     SettlementAdjustmentKind = "REFUND"
	SettlementAdjustmentKindFee        SettlementAdjustmentKind = "FEE"
	SettlementAdjustmentKindReversal   SettlementAdjustmentKind = "REVERSAL"
	SettlementAdjustmentKindCorrection SettlementAdjustmentKind = "CORRECTION"
)

type SettlementAdjustmentRecoveryStatus string

const (
	SettlementAdjustmentRecoveryStatusNone       SettlementAdjustmentRecoveryStatus = "NONE"
	SettlementAdjustmentRecoveryStatusPending    SettlementAdjustmentRecoveryStatus = "PENDING"
	SettlementAdjustmentRecoveryStatusProcessing SettlementAdjustmentRecoveryStatus = "PROCESSING"
	SettlementAdjustmentRecoveryStatusRecovered  SettlementAdjustmentRecoveryStatus = "RECOVERED"
	SettlementAdjustmentRecoveryStatusFailed     SettlementAdjustmentRecoveryStatus = "FAILED"
)

type SettlementBatchStatus string

const (
	SettlementBatchStatusCreated    SettlementBatchStatus = "CREATED"
	SettlementBatchStatusSubmitted  SettlementBatchStatus = "SUBMITTED"
	SettlementBatchStatusProcessing SettlementBatchStatus = "PROCESSING"
	SettlementBatchStatusPartial    SettlementBatchStatus = "PARTIAL"
	SettlementBatchStatusSettled    SettlementBatchStatus = "SETTLED"
	SettlementBatchStatusFailed     SettlementBatchStatus = "FAILED"
	SettlementBatchStatusCancelled  SettlementBatchStatus = "CANCELLED"
)

type SettlementBatchItemStatus string

const (
	SettlementBatchItemStatusCreated    SettlementBatchItemStatus = "CREATED"
	SettlementBatchItemStatusSubmitted  SettlementBatchItemStatus = "SUBMITTED"
	SettlementBatchItemStatusProcessing SettlementBatchItemStatus = "PROCESSING"
	SettlementBatchItemStatusSuccess    SettlementBatchItemStatus = "SUCCESS"
	SettlementBatchItemStatusFailed     SettlementBatchItemStatus = "FAILED"
	SettlementBatchItemStatusCancelled  SettlementBatchItemStatus = "CANCELLED"
)

type ReconciliationRunBillType string

const (
	ReconciliationRunBillTypePayment    ReconciliationRunBillType = "PAYMENT"
	ReconciliationRunBillTypeRefund     ReconciliationRunBillType = "REFUND"
	ReconciliationRunBillTypeSettlement ReconciliationRunBillType = "SETTLEMENT"
)

type ReconciliationRunStatus string

const (
	ReconciliationRunStatusPending    ReconciliationRunStatus = "PENDING"
	ReconciliationRunStatusDownloaded ReconciliationRunStatus = "DOWNLOADED"
	ReconciliationRunStatusParsed     ReconciliationRunStatus = "PARSED"
	ReconciliationRunStatusCompleted  ReconciliationRunStatus = "COMPLETED"
	ReconciliationRunStatusFailed     ReconciliationRunStatus = "FAILED"
)

type ReconciliationDiffLocalType string

const (
	ReconciliationDiffLocalTypePayment              ReconciliationDiffLocalType = "PAYMENT"
	ReconciliationDiffLocalTypeRefund               ReconciliationDiffLocalType = "REFUND"
	ReconciliationDiffLocalTypeSettlement           ReconciliationDiffLocalType = "SETTLEMENT"
	ReconciliationDiffLocalTypeSettlementSupplement ReconciliationDiffLocalType = "SETTLEMENT_SUPPLEMENT"
)

type ReconciliationDiffDiffType string

const (
	ReconciliationDiffDiffTypeLocalOnly   ReconciliationDiffDiffType = "LOCAL_ONLY"
	ReconciliationDiffDiffTypeChannelOnly ReconciliationDiffDiffType = "CHANNEL_ONLY"
	ReconciliationDiffDiffTypeAmount      ReconciliationDiffDiffType = "AMOUNT"
	ReconciliationDiffDiffTypeStatus      ReconciliationDiffDiffType = "STATUS"
	ReconciliationDiffDiffTypeFee         ReconciliationDiffDiffType = "FEE"
)

type ReconciliationDiffStatus string

const (
	ReconciliationDiffStatusOpen          ReconciliationDiffStatus = "OPEN"
	ReconciliationDiffStatusInvestigating ReconciliationDiffStatus = "INVESTIGATING"
	ReconciliationDiffStatusResolved      ReconciliationDiffStatus = "RESOLVED"
	ReconciliationDiffStatusIgnored       ReconciliationDiffStatus = "IGNORED"
)

type ActivityStatus string

const (
	ActivityStatusDraft     ActivityStatus = "DRAFT"
	ActivityStatusPublished ActivityStatus = "PUBLISHED"
	ActivityStatusOffline   ActivityStatus = "OFFLINE"
)

type BannerTargetType string

const (
	BannerTargetTypeGame     BannerTargetType = "GAME"
	BannerTargetTypeCoach    BannerTargetType = "COACH"
	BannerTargetTypeService  BannerTargetType = "SERVICE"
	BannerTargetTypeActivity BannerTargetType = "ACTIVITY"
)

type BannerStatus string

const (
	BannerStatusDraft     BannerStatus = "DRAFT"
	BannerStatusPublished BannerStatus = "PUBLISHED"
	BannerStatusOffline   BannerStatus = "OFFLINE"
)

type CoachRecommendationStatus string

const (
	CoachRecommendationStatusEnabled  CoachRecommendationStatus = "ENABLED"
	CoachRecommendationStatusDisabled CoachRecommendationStatus = "DISABLED"
)

type ContentAuditSubjectType string

const (
	ContentAuditSubjectTypeUser          ContentAuditSubjectType = "USER"
	ContentAuditSubjectTypeCoach         ContentAuditSubjectType = "COACH"
	ContentAuditSubjectTypeService       ContentAuditSubjectType = "SERVICE"
	ContentAuditSubjectTypeSkill         ContentAuditSubjectType = "SKILL"
	ContentAuditSubjectTypeReviewContent ContentAuditSubjectType = "REVIEW_CONTENT"
	ContentAuditSubjectTypeAsset         ContentAuditSubjectType = "ASSET"
	ContentAuditSubjectTypeActivity      ContentAuditSubjectType = "ACTIVITY"
)

type ContentAuditStatus string

const (
	ContentAuditStatusPending      ContentAuditStatus = "PENDING"
	ContentAuditStatusAutoPassed   ContentAuditStatus = "AUTO_PASSED"
	ContentAuditStatusManualReview ContentAuditStatus = "MANUAL_REVIEW"
	ContentAuditStatusApproved     ContentAuditStatus = "APPROVED"
	ContentAuditStatusRejected     ContentAuditStatus = "REJECTED"
	ContentAuditStatusSuperseded   ContentAuditStatus = "SUPERSEDED"
)

type UserPenaltyAction string

const (
	UserPenaltyActionBanUser      UserPenaltyAction = "BAN_USER"
	UserPenaltyActionSuspendCoach UserPenaltyAction = "SUSPEND_COACH"
	UserPenaltyActionBlockService UserPenaltyAction = "BLOCK_SERVICE"
)

type UserPenaltyStatus string

const (
	UserPenaltyStatusActive  UserPenaltyStatus = "ACTIVE"
	UserPenaltyStatusRevoked UserPenaltyStatus = "REVOKED"
	UserPenaltyStatusExpired UserPenaltyStatus = "EXPIRED"
)

type AdminUserStatus string

const (
	AdminUserStatusActive   AdminUserStatus = "ACTIVE"
	AdminUserStatusDisabled AdminUserStatus = "DISABLED"
)

type AdminRoleDataScope string

const (
	AdminRoleDataScopeAll      AdminRoleDataScope = "ALL"
	AdminRoleDataScopeAssigned AdminRoleDataScope = "ASSIGNED"
)

type AdminRoleStatus string

const (
	AdminRoleStatusActive   AdminRoleStatus = "ACTIVE"
	AdminRoleStatusDisabled AdminRoleStatus = "DISABLED"
)

type AdminApprovalStatus string

const (
	AdminApprovalStatusPending  AdminApprovalStatus = "PENDING"
	AdminApprovalStatusApproved AdminApprovalStatus = "APPROVED"
	AdminApprovalStatusRejected AdminApprovalStatus = "REJECTED"
	AdminApprovalStatusExpired  AdminApprovalStatus = "EXPIRED"
	AdminApprovalStatusExecuted AdminApprovalStatus = "EXECUTED"
)

type AdminAuditLogActorType string

const (
	AdminAuditLogActorTypeAdmin  AdminAuditLogActorType = "ADMIN"
	AdminAuditLogActorTypeSystem AdminAuditLogActorType = "SYSTEM"
)

type AdminAuditLogOutcome string

const (
	AdminAuditLogOutcomeSuccess AdminAuditLogOutcome = "SUCCESS"
	AdminAuditLogOutcomeFailed  AdminAuditLogOutcome = "FAILED"
)

type BusinessRuleVersionCode string

const (
	BusinessRuleVersionCodeOrder        BusinessRuleVersionCode = "ORDER"
	BusinessRuleVersionCodeCancellation BusinessRuleVersionCode = "CANCELLATION"
	BusinessRuleVersionCodeSettlement   BusinessRuleVersionCode = "SETTLEMENT"
	BusinessRuleVersionCodeSupport      BusinessRuleVersionCode = "SUPPORT"
)

type BusinessRuleVersionStatus string

const (
	BusinessRuleVersionStatusDraft     BusinessRuleVersionStatus = "DRAFT"
	BusinessRuleVersionStatusPublished BusinessRuleVersionStatus = "PUBLISHED"
	BusinessRuleVersionStatusRetired   BusinessRuleVersionStatus = "RETIRED"
)

type IdempotencyKeyActorType string

const (
	IdempotencyKeyActorTypeUser   IdempotencyKeyActorType = "USER"
	IdempotencyKeyActorTypeAdmin  IdempotencyKeyActorType = "ADMIN"
	IdempotencyKeyActorTypeSystem IdempotencyKeyActorType = "SYSTEM"
)

type IdempotencyKeyStatus string

const (
	IdempotencyKeyStatusProcessing IdempotencyKeyStatus = "PROCESSING"
	IdempotencyKeyStatusSucceeded  IdempotencyKeyStatus = "SUCCEEDED"
	IdempotencyKeyStatusFailed     IdempotencyKeyStatus = "FAILED"
)

type OutboxEventAggregateType string

const (
	OutboxEventAggregateTypeOrder                OutboxEventAggregateType = "ORDER"
	OutboxEventAggregateTypePayment              OutboxEventAggregateType = "PAYMENT"
	OutboxEventAggregateTypeRefund               OutboxEventAggregateType = "REFUND"
	OutboxEventAggregateTypeSettlement           OutboxEventAggregateType = "SETTLEMENT"
	OutboxEventAggregateTypeSettlementSupplement OutboxEventAggregateType = "SETTLEMENT_SUPPLEMENT"
	OutboxEventAggregateTypeCoachApplication     OutboxEventAggregateType = "COACH_APPLICATION"
	OutboxEventAggregateTypeIdentityCheck        OutboxEventAggregateType = "IDENTITY_CHECK"
	OutboxEventAggregateTypeCoachSkill           OutboxEventAggregateType = "COACH_SKILL"
	OutboxEventAggregateTypeAfterSalesCase       OutboxEventAggregateType = "AFTER_SALES_CASE"
	OutboxEventAggregateTypeContentAudit         OutboxEventAggregateType = "CONTENT_AUDIT"
	OutboxEventAggregateTypeNotification         OutboxEventAggregateType = "NOTIFICATION"
	OutboxEventAggregateTypeReconciliation       OutboxEventAggregateType = "RECONCILIATION"
	OutboxEventAggregateTypeUser                 OutboxEventAggregateType = "USER"
	OutboxEventAggregateTypeCoach                OutboxEventAggregateType = "COACH"
	OutboxEventAggregateTypeService              OutboxEventAggregateType = "SERVICE"
	OutboxEventAggregateTypeBusinessRule         OutboxEventAggregateType = "BUSINESS_RULE"
)

type OutboxEventStatus string

const (
	OutboxEventStatusPending    OutboxEventStatus = "PENDING"
	OutboxEventStatusProcessing OutboxEventStatus = "PROCESSING"
	OutboxEventStatusDone       OutboxEventStatus = "DONE"
	OutboxEventStatusDead       OutboxEventStatus = "DEAD"
)

type UserActivityEventEventType string

const (
	UserActivityEventEventTypeAppActive   UserActivityEventEventType = "APP_ACTIVE"
	UserActivityEventEventTypeViewService UserActivityEventEventType = "VIEW_SERVICE"
	UserActivityEventEventTypeCreateOrder UserActivityEventEventType = "CREATE_ORDER"
)

type DailyMetricMetricCode string

const (
	DailyMetricMetricCodeDau         DailyMetricMetricCode = "DAU"
	DailyMetricMetricCodePaidGmv     DailyMetricMetricCode = "PAID_GMV"
	DailyMetricMetricCodeRefund      DailyMetricMetricCode = "REFUND"
	DailyMetricMetricCodeCommission  DailyMetricMetricCode = "COMMISSION"
	DailyMetricMetricCodeCoachActive DailyMetricMetricCode = "COACH_ACTIVE"
)

// SettlementSupplementStatus 独立补款指令状态；未知结果保留处理中，重试复用 RequestNo。
type SettlementSupplementStatus string

const (
	SettlementSupplementStatusPending    SettlementSupplementStatus = "PENDING"
	SettlementSupplementStatusSubmitted  SettlementSupplementStatus = "SUBMITTED"
	SettlementSupplementStatusProcessing SettlementSupplementStatus = "PROCESSING"
	SettlementSupplementStatusSuccess    SettlementSupplementStatus = "SUCCESS"
	SettlementSupplementStatusFailed     SettlementSupplementStatus = "FAILED"
	SettlementSupplementStatusCancelled  SettlementSupplementStatus = "CANCELLED"
)
